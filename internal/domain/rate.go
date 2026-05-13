package domain

import (
	"context"
	"time"
)

type Rate struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Amount    int64     `json:"amount" gorm:"not null"`
	Discount  *int64    `json:"discount"`
	Notes     string    `json:"notes" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateRateRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Amount   int64  `json:"amount" validate:"required,gt=0"`
	Discount *int64 `json:"discount" validate:"omitempty,gte=0"`
	Notes    string `json:"notes" validate:"omitempty,max=1000"`
}

type UpdateRateRequest struct {
	Name     string `json:"name" validate:"omitempty,max=255"`
	Amount   int64  `json:"amount" validate:"omitempty,gt=0"`
	Discount *int64 `json:"discount" validate:"omitempty,gte=0"`
	Notes    string `json:"notes" validate:"omitempty,max=1000"`
}

type RateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Amount    int64     `json:"amount"`
	Discount  *int64    `json:"discount,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RateRepository interface {
	Create(ctx context.Context, r *Rate) error
	FindByID(ctx context.Context, id uint) (*Rate, error)
	FindAll(ctx context.Context, page, limit int) ([]Rate, int64, error)
	Update(ctx context.Context, r *Rate) error
	Delete(ctx context.Context, id uint) error
}

type RateService interface {
	Create(ctx context.Context, req CreateRateRequest) (*RateResponse, error)
	FindByID(ctx context.Context, id uint) (*RateResponse, error)
	FindAll(ctx context.Context, page, limit int) ([]RateResponse, int64, error)
	Update(ctx context.Context, id uint, req UpdateRateRequest) (*RateResponse, error)
	Delete(ctx context.Context, id uint) error
}
