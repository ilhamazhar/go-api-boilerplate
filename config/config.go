package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort          string
	DatabaseURL         string
	RedisURL            string
	JWTSecret           string
	JWTRefreshSecret    string
	JWTExpiry           time.Duration
	JWTRefreshExpiry    time.Duration
	XenditAPIKey        string
	XenditCallbackToken string
	XenditWebhookToken  string
}

func Load() (Config, error) {
	jwtExpiry, err := strconv.ParseFloat(getEnv("JWT_EXPIRY_HOURS", "24"), 64)
	if err != nil {
		return Config{}, errors.New("invalid JWT_EXPIRY_HOURS")
	}
	jwtRefreshExpiry, err := strconv.ParseFloat(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"), 64)
	if err != nil {
		return Config{}, errors.New("invalid JWT_REFRESH_EXPIRY_HOURS")
	}

	cfg := Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		RedisURL:            getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTRefreshSecret:    os.Getenv("JWT_REFRESH_SECRET"),
		JWTExpiry:           time.Duration(jwtExpiry * float64(time.Hour)),
		JWTRefreshExpiry:    time.Duration(jwtRefreshExpiry * float64(time.Hour)),
		XenditAPIKey:        os.Getenv("XENDIT_API_KEY"),
		XenditCallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
		XenditWebhookToken:  os.Getenv("XENDIT_WEBHOOK_TOKEN"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}
	if cfg.JWTRefreshSecret == "" {
		return Config{}, errors.New("JWT_REFRESH_SECRET is required")
	}
	if cfg.XenditAPIKey == "" {
		return Config{}, errors.New("XENDIT_API_KEY is required")
	}
	if cfg.XenditCallbackToken == "" {
		return Config{}, errors.New("XENDIT_CALLBACK_TOKEN is required")
	}
	if cfg.XenditWebhookToken == "" {
		return Config{}, errors.New("XENDIT_WEBHOOK_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
