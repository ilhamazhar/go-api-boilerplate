package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ilhamazhar/golang-gpt/internal/domain"
)

type rateService struct {
	repo domain.RateRepository
}

func NewRateService(repo domain.RateRepository) domain.RateService {
	return &rateService{repo: repo}
}

func (s *rateService) Create(ctx context.Context, req domain.CreateRateRequest) (*domain.RateResponse, error) {
	r := &domain.Rate{
		Name:     req.Name,
		Amount:   req.Amount,
		Discount: req.Discount,
		Notes:    req.Notes,
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to create rate: %w", err)
	}
	return toRateResponse(r), nil
}

func (s *rateService) FindAll(ctx context.Context, page, limit int) ([]domain.RateResponse, int64, error) {
	rates, total, err := s.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch rates: %w", err)
	}
	result := make([]domain.RateResponse, len(rates))
	for i, r := range rates {
		result[i] = *toRateResponse(&r)
	}
	return result, total, nil
}

func (s *rateService) FindByID(ctx context.Context, id uint) (*domain.RateResponse, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("rate not found")
	}
	return toRateResponse(r), nil
}

func (s *rateService) Update(ctx context.Context, id uint, req domain.UpdateRateRequest) (*domain.RateResponse, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("rate not found")
	}
	if req.Name != "" {
		r.Name = req.Name
	}
	if req.Amount > 0 {
		r.Amount = req.Amount
	}
	r.Discount = req.Discount
	if req.Notes != "" {
		r.Notes = req.Notes
	}
	if err := s.repo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("failed to update rate: %w", err)
	}
	return toRateResponse(r), nil
}

func (s *rateService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return errors.New("rate not found")
	}
	return s.repo.Delete(ctx, id)
}

func toRateResponse(r *domain.Rate) *domain.RateResponse {
	return &domain.RateResponse{
		ID:        r.ID,
		Name:      r.Name,
		Amount:    r.Amount,
		Discount:  r.Discount,
		Notes:     r.Notes,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
