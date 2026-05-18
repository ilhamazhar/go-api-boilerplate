package repository

import (
	"context"

	"github.com/ilhamazhar/golang-gpt/internal/domain"
	"gorm.io/gorm"
)

type rateRepo struct {
	db *gorm.DB
}

func NewRateRepository(db *gorm.DB) domain.RateRepository {
	return &rateRepo{db: db}
}

func (r *rateRepo) Create(ctx context.Context, rate *domain.Rate) error {
	return r.db.WithContext(ctx).Create(rate).Error
}

func (r *rateRepo) FindAll(ctx context.Context, page, limit int) ([]domain.Rate, int64, error) {
	var rates []domain.Rate
	var total int64

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Model(&domain.Rate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&rates).Error
	return rates, total, err
}

func (r *rateRepo) FindByID(ctx context.Context, id uint) (*domain.Rate, error) {
	var rate domain.Rate
	err := r.db.WithContext(ctx).First(&rate, id).Error
	return &rate, err
}

func (r *rateRepo) Update(ctx context.Context, rate *domain.Rate) error {
	return r.db.WithContext(ctx).Save(rate).Error
}

func (r *rateRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Rate{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
