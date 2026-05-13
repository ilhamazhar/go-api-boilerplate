package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"github.com/ilhamazhar/golang-gpt/internal/middleware"
	"github.com/ilhamazhar/golang-gpt/pkg/response"
	"github.com/ilhamazhar/golang-gpt/pkg/validator"
)

type PaymentHandler struct {
	svc domain.PaymentService
}

func NewPaymentHandler(svc domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) CreateQRIS(c *gin.Context) {
	claims := middleware.ClaimsFromContext(c)

	var req domain.CreateQRISRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid JSON", nil)
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.Fail(c, http.StatusUnprocessableEntity, "validation failed", errs)
		return
	}

	result, err := h.svc.CreateQRIS(c.Request.Context(), claims.UserID, req)
	if err != nil {
		response.Fail(c, http.StatusBadGateway, err.Error(), nil)
		return
	}
	response.OK(c, http.StatusCreated, "QRIS created", result)
}

func (h *PaymentHandler) GetStatus(c *gin.Context) {
	orderRef := c.Param("order_ref")
	result, err := h.svc.GetStatus(c.Request.Context(), orderRef)
	if err != nil {
		response.Fail(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	response.OK(c, http.StatusOK, "Order status retrieved", result)
}

// Webhook receives Xendit's callback. Public endpoint, authed via X-Callback-Token header.
func (h *PaymentHandler) Webhook(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.svc.HandleWebhook(c.Request.Context(), callbackToken, body); err != nil {
		// Log internally but always return 200 so Xendit doesn't retry forever
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusOK)
}
