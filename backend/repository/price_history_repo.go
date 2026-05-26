package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type PriceHistoryRepository interface {
	Create(h *model.PriceHistory) error
	FindByItem(itemID uint) ([]model.PriceHistory, error)
}

type priceHistoryRepo struct {
	db *gorm.DB
}

func NewPriceHistoryRepository(db *gorm.DB) PriceHistoryRepository {
	return &priceHistoryRepo{db: db}
}

func (r *priceHistoryRepo) Create(h *model.PriceHistory) error {
	return r.db.Create(h).Error
}

func (r *priceHistoryRepo) FindByItem(itemID uint) ([]model.PriceHistory, error) {
	var records []model.PriceHistory
	err := r.db.Where("item_id = ?", itemID).
		Order("recorded_at DESC").
		Find(&records).Error
	return records, err
}
