package main

import (
	"log"

	"github.com/ilhamazhar/golang-gpt/config"
	"github.com/ilhamazhar/golang-gpt/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	log.Printf("Config loaded: port=%s jwt_expiry=%v jwt_refresh_expiry=%v",
		cfg.ServerPort, cfg.JWTExpiry, cfg.JWTRefreshExpiry)

	a, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	log.Fatal(a.Run())
}
