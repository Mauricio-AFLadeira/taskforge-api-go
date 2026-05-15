package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds runtime settings loaded from the environment.
type Config struct {
	AppEnv  string
	AppPort int

	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret         string
	JWTAccessTokenTTL string
	RefreshTokenTTL   string
}

// Load reads configuration from environment variables.
func Load() (Config, error) {
	appPort := 8080
	if v := os.Getenv("APP_PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("APP_PORT: %w", err)
		}
		appPort = n
	}

	redisDB := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("REDIS_DB: %w", err)
		}
		redisDB = n
	}

	cfg := Config{
		AppEnv:            getenvDefault("APP_ENV", "development"),
		AppPort:           appPort,
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		RedisAddr:         getenvDefault("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		RedisDB:           redisDB,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTAccessTokenTTL: getenvDefault("JWT_ACCESS_TOKEN_TTL", "15m"),
		RefreshTokenTTL:   getenvDefault("REFRESH_TOKEN_TTL", "168h"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
