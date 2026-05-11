package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

func Load() Config {
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))

	return Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", "postgresql://postgres:mysecretpassword@localhost:5432/learning"),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		JWTExpiry:           time.Duration(jwtExpiry) * time.Hour,
		JWTRefreshExpiry:    time.Duration(jwtRefreshExpiry) * time.Hour,
		XenditAPIKey:        getEnv("XENDIT_API_KEY", "xnd_development_r9gg0jAMnyiDThl6YrZao6TMEawLaJj471qDrXrBGWfs5aBQMkyLf0YC7fyENIxW"),
		XenditCallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", "EqGQ7LX6coXZbpD7Bhdlw2uJ8ITMMGim5Lm1kkW84MnqmQVm"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
