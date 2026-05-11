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
	xenclient "github.com/ilhamazhar/golang-gpt/pkg/xendit"
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

	db.AutoMigrate(&domain.User{}, &domain.Payment{})

	 // --- External clients ---
	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiry)
	refreshManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTRefreshExpiry)
	xenditClient := xenclient.NewClient(cfg.XenditAPIKey, cfg.XenditCallbackToken)

	// --- Repositories ---
	userRepo := repository.NewUserRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// --- Services ---
	authService := service.NewAuthService(userRepo, jwtManager, refreshManager)
	paymentService := service.NewPaymentService(paymentRepo, xenditClient)

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()
	r.Use(corsMiddleware())
	registerRoutes(r, Handlers{Auth: authHandler, Payment: paymentHandler}, jwtManager)

	return &App{cfg: cfg, router: r}, nil
}

func (a *App) Run() error {
	return a.router.Run(":" + a.cfg.ServerPort)
}

func initDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
