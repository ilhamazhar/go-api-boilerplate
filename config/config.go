package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort       string
	DatabaseURL      string
	JWTSecret        string
	JWTExpiry        time.Duration
	JWTRefreshExpiry time.Duration
}

func Load() Config {
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY", "168"))

	return Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgresql://postgres:mysecretpassword@localhost:5432/learning"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
		JWTExpiry:        time.Duration(jwtExpiry) * time.Hour,
		JWTRefreshExpiry: time.Duration(jwtRefreshExpiry) * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
