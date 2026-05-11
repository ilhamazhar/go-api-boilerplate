package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ilhamazhar/golang-gpt/internal/handler"
	"github.com/ilhamazhar/golang-gpt/internal/middleware"
	"github.com/ilhamazhar/golang-gpt/pkg/jwt"
)

type Handlers struct {
	Auth    *handler.AuthHandler
	Payment *handler.PaymentHandler
}

func registerRoutes(r *gin.Engine, h Handlers, jwtManager *jwt.Manager) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.Auth(jwtManager))
	{
		api.GET("/me", h.Auth.Me)
		api.POST("/payments/qris", h.Payment.CreateQRIS)
		api.GET("/payments/:order_ref", h.Payment.GetStatus)
	}

	r.POST("/webhooks/xendit", h.Payment.Webhook)
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
