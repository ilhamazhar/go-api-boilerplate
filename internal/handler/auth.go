package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"github.com/ilhamazhar/golang-gpt/internal/middleware"
	"github.com/ilhamazhar/golang-gpt/pkg/response"
	"github.com/ilhamazhar/golang-gpt/pkg/validator"
)

type AuthHandler struct {
	auth domain.AuthService
}

func NewAuthHandler(auth domain.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	if errs := validator.Validate(req); errs != nil {
		response.Fail(c, http.StatusUnprocessableEntity, "Validation failed", errs)
		return
	}

	user, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusConflict, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusCreated, "Registered successfully", domain.ToUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	if errs := validator.Validate(req); errs != nil {
		response.Fail(c, http.StatusUnprocessableEntity, "Validation failed", errs)
		return
	}

	tokens, err := h.auth.Login(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusOK, "Logged in successfully", tokens)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)
	response.OK(c, http.StatusOK, "User info retrieved", domain.UserResponse{
		ID:    claims.UserID,
		Email: claims.Email,
	})
}
