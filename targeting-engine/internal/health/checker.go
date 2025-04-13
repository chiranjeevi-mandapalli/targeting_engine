package health

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type DatabaseChecker struct {
	db *sql.DB
}

func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

func (c *DatabaseChecker) Check(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

type RedisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{client: client}
}

func (c *RedisChecker) Check(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	return err
}

type HealthService struct {
	checkers map[string]Checker
}

func NewHealthService(db *sql.DB, redisClient *redis.Client) *HealthService {
	return &HealthService{
		checkers: map[string]Checker{
			"database": NewDatabaseChecker(db),
			"redis":    NewRedisChecker(redisClient),
		},
	}
}

func (s *HealthService) Check(ctx context.Context) HealthResponse {
	response := HealthResponse{
		Status:  StatusUp,
		Details: make(map[string]Status),
	}

	for name, checker := range s.checkers {
		if err := checker.Check(ctx); err != nil {
			response.Status = StatusDown
			response.Details[name] = StatusDown
		} else {
			response.Details[name] = StatusUp
		}
	}

	return response
}
