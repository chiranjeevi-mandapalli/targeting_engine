package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"targeting-engine/internal/campaign"
	"targeting-engine/internal/delivery"
	"targeting-engine/internal/health"
	"targeting-engine/internal/monitoring"
	"targeting-engine/internal/seed"
	"targeting-engine/internal/targeting"
	"targeting-engine/pkg/config"
	"targeting-engine/pkg/logging"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize logger
	logger := logging.New("targeting-engine")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database connection
	db, err := initDatabase(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Seed database if enabled
	if cfg.Database.Seed {
		seedDatabase(db, logger)
	}

	// Initialize monitoring
	metrics := monitoring.Init()
	logger.Info().Msg("Monitoring initialized")

	// Initialize Redis
	redisClient := initRedis(cfg, logger)
	defer redisClient.Close()

	// Initialize health service
	healthService := health.NewHealthService(db.DB, redisClient)
	logger.Info().Msg("Health checks initialized")

	// Initialize services
	campaignSvc, targetingSvc := initServices(db, redisClient, cfg, logger, metrics)

	// Create HTTP server
	server := initHTTPServer(cfg, campaignSvc, targetingSvc, healthService, logger, metrics)

	// Start server in goroutine
	go func() {
		logger.Info().
			Str("host", cfg.Server.Host).
			Int("port", cfg.Server.Port).
			Msg("Starting server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	waitForShutdown(server, logger)
}

func initDatabase(cfg *config.Config, logger *logging.Logger) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info().Msg("Database connection established")
	return db, nil
}

func seedDatabase(db *sqlx.DB, logger *logging.Logger) {
	logger.Info().Msg("Starting database seeding...")

	startTime := time.Now()
	seeder := seed.NewSeeder(db.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := seeder.SeedAll(ctx); err != nil {
		logger.Error().Err(err).Msg("Database seeding failed")
		return
	}

	logger.Info().
		Str("duration", time.Since(startTime).String()).
		Msg("Database seeded successfully")
}

func initRedis(cfg *config.Config, logger *logging.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	logger.Info().Msg("Redis connection established")
	return rdb
}

func initServices(db *sqlx.DB, rdb *redis.Client, cfg *config.Config, logger *logging.Logger, metrics *monitoring.Metrics) (*campaign.Service, *targeting.Evaluator) {
	// Initialize repositories
	campaignRepo := campaign.NewPostgresRepository(db)
	ruleRepo := targeting.NewPostgresRuleRepository(db)

	// Create cached repositories with Redis client directly
	cachedCampaignRepo := campaign.NewCachedRepository(campaignRepo, rdb, 5*time.Minute)
	cachedRuleRepo := targeting.NewCachedRuleRepository(ruleRepo, rdb, 10*time.Minute)

	// Initialize services
	campaignSvc := campaign.NewService(cachedCampaignRepo)
	targetingSvc := targeting.NewEvaluator(cachedRuleRepo)

	logger.Info().Msg("Services initialized")
	return campaignSvc, targetingSvc
}

func initHTTPServer(cfg *config.Config, campaignSvc *campaign.Service, targetingSvc *targeting.Evaluator, healthService *health.HealthService, logger *logging.Logger, metrics *monitoring.Metrics) *http.Server {
	// Create delivery service
	deliverySvc := delivery.NewService(campaignSvc, targetingSvc)

	// Create router
	router := mux.NewRouter()

	// Add middleware
	router.Use(
		loggingMiddleware(logger),
		monitoring.MetricsMiddleware(metrics),
	)

	// Register handlers
	router.PathPrefix("/v1/campaigns").Handler(campaign.MakeHTTPHandler(campaignSvc))
	router.Handle("/v1/delivery", delivery.MakeHTTPHandler(deliverySvc))

	// Add monitoring endpoints
	router.Handle("/metrics", promhttp.Handler())

	// Add health endpoints
	router.Handle("/healthz", health.MakeHandler(healthService))

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func loggingMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("duration", time.Since(start).String()).
				Msg("Request processed")
		})
	}
}

func waitForShutdown(server *http.Server, logger *logging.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown failed")
	}

	logger.Info().Msg("Server exited properly")
}
