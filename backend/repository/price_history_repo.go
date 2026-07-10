package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type PriceHistoryRepository interface {
	Create(h *model.PriceHistory) error
	FindByItem(itemID uint) ([]model.PriceHistory, error)
	FindByID(id uint) (*model.PriceHistory, error)
	FindVisibleByItemAndDate(itemID uint, date string) ([]model.PriceHistory, error)
	Delete(id uint) error
	ToggleHidden(id uint, isHidden bool) error
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
	err := r.db.Where("item_id = ? AND is_hidden = false", itemID).
		Order("recorded_at DESC").
		Find(&records).Error
	return records, err
}

func (r *priceHistoryRepo) FindByID(id uint) (*model.PriceHistory, error) {
	var h model.PriceHistory
	err := r.db.First(&h, id).Error
	return &h, err
}

func (r *priceHistoryRepo) FindVisibleByItemAndDate(itemID uint, date string) ([]model.PriceHistory, error) {
	var records []model.PriceHistory
	err := r.db.Where(
		"item_id = ? AND is_hidden = false AND (recorded_at AT TIME ZONE 'Asia/Taipei')::date = ?",
		itemID, date,
	).Order("recorded_at DESC").Find(&records).Error
	return records, err
}

func (r *priceHistoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.PriceHistory{}, id).Error
}

func (r *priceHistoryRepo) ToggleHidden(id uint, isHidden bool) error {
	return r.db.Model(&model.PriceHistory{}).Where("id = ?", id).Update("is_hidden", isHidden).Error
}
