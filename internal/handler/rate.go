package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"github.com/ilhamazhar/golang-gpt/pkg/response"
)

type RateHandler struct {
	svc domain.RateService
}

func NewRateHandler(svc domain.RateService) *RateHandler {
	return &RateHandler{svc: svc}
}

func (h *RateHandler) Create(c *gin.Context) {
	var req domain.CreateRateRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusCreated, "Rate created", result)
}

func (h *RateHandler) GetAll(c *gin.Context) {
	page, limit := parsePagination(c)

	result, total, err := h.svc.FindAll(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	pagination := response.NewPagination(page, limit, total)
	response.OKPaginated(c, http.StatusOK, "Rates retrieved", result, pagination)
}

func (h *RateHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id", nil)
		return
	}

	result, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusOK, "Rate retrieved", result)
}

func (h *RateHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id", nil)
		return
	}

	var req domain.UpdateRateRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusOK, "Rate updated", result)
}

func (h *RateHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid id", nil)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.OK(c, http.StatusOK, "Rate deleted", nil)
}
