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
	JWTSecret           string
	JWTExpiry           time.Duration
	JWTRefreshExpiry    time.Duration
	XenditAPIKey        string
	XenditCallbackToken string
}

func Load() (Config, error) {
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))

	cfg := Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTExpiry:           time.Duration(jwtExpiry) * time.Hour,
		JWTRefreshExpiry:    time.Duration(jwtRefreshExpiry) * time.Hour,
		XenditAPIKey:        os.Getenv("XENDIT_API_KEY"),
		XenditCallbackToken: os.Getenv("XENDIT_CALLBACK_TOKEN"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}
	if cfg.XenditAPIKey == "" {
		return Config{}, errors.New("XENDIT_API_KEY is required")
	}
	if cfg.XenditCallbackToken == "" {
		return Config{}, errors.New("XENDIT_CALLBACK_TOKEN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
