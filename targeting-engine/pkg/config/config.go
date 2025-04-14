package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server struct {
		Host         string
		Port         int
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
		Seed     bool
	}

	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}

func Load() (*Config, error) {
	var cfg Config

	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", 8080)
	cfg.Server.ReadTimeout = time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 5)) * time.Second
	cfg.Server.WriteTimeout = time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second
	cfg.Server.IdleTimeout = time.Duration(getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "9063770754")
	cfg.Database.Name = getEnv("DB_NAME", "targeting")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	cfg.Redis.Host = getEnv("REDIS_HOST", "localhost")
	cfg.Redis.Port = getEnvAsInt("REDIS_PORT", 6379)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	return &cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}
