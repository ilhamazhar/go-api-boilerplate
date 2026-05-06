package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ilhamazhar/golang-gpt/config"
	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"github.com/ilhamazhar/golang-gpt/internal/handler"
	"github.com/ilhamazhar/golang-gpt/internal/repository"
	"github.com/ilhamazhar/golang-gpt/internal/service"
	"github.com/ilhamazhar/golang-gpt/pkg/jwt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	cfg    config.Config
	router *gin.Engine
}

func New(cfg config.Config) (*App, error) {
	db, err := initDB(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	db.AutoMigrate(&domain.User{})

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	refreshManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTRefreshExpiry)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager, refreshManager)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(corsMiddleware())
	registerRoutes(r, authHandler, jwtManager)

	return &App{cfg: cfg, router: r}, nil
}

func (a *App) Run() error {
	return a.router.Run(":" + a.cfg.ServerPort)
}

func initDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
