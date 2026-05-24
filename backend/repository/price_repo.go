package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type PriceRepository interface {
	FindByItemAndDate(itemID uint, date string) (*model.PriceRecord, error)
	FindByItemsAndDates(itemIDs []uint, dates []string) ([]model.PriceRecord, error)
	Create(record *model.PriceRecord) error
	UpdatePrice(record *model.PriceRecord, price float64) error
	FindHistoryByItem(itemID uint, limit int) ([]model.PriceRecord, error)
	DeleteByItem(itemID uint) error
}

type priceRepo struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepo{db: db}
}

func (r *priceRepo) FindByItemAndDate(itemID uint, date string) (*model.PriceRecord, error) {
	var record model.PriceRecord
	err := r.db.Where("item_id = ? AND recorded_date = ?", itemID, date).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *priceRepo) FindByItemsAndDates(itemIDs []uint, dates []string) ([]model.PriceRecord, error) {
	var records []model.PriceRecord
	err := r.db.Where("item_id IN ? AND recorded_date IN ?", itemIDs, dates).Find(&records).Error
	return records, err
}

func (r *priceRepo) Create(record *model.PriceRecord) error {
	return r.db.Create(record).Error
}

func (r *priceRepo) UpdatePrice(record *model.PriceRecord, price float64) error {
	return r.db.Model(record).Update("price", price).Error
}

func (r *priceRepo) FindHistoryByItem(itemID uint, limit int) ([]model.PriceRecord, error) {
	var records []model.PriceRecord
	err := r.db.Where("item_id = ?", itemID).
		Order("recorded_date desc").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (r *priceRepo) DeleteByItem(itemID uint) error {
	return r.db.Where("item_id = ?", itemID).Delete(&model.PriceRecord{}).Error
}
