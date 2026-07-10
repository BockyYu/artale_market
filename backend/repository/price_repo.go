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
	FindLatestByItem(itemID uint) (*model.PriceRecord, error)
	FindLatestByItems(itemIDs []uint) ([]model.PriceRecord, error)
	DeleteByItem(itemID uint) error
	DeleteByItemAndDate(itemID uint, date string) error
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

func (r *priceRepo) FindLatestByItem(itemID uint) (*model.PriceRecord, error) {
	var record model.PriceRecord
	err := r.db.Where("item_id = ?", itemID).Order("recorded_date desc").First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *priceRepo) FindLatestByItems(itemIDs []uint) ([]model.PriceRecord, error) {
	if len(itemIDs) == 0 {
		return nil, nil
	}
	var records []model.PriceRecord
	// 每個 item_id 只取 recorded_date 最新的那筆
	err := r.db.
		Where("(item_id, recorded_date) IN (?)",
			r.db.Model(&model.PriceRecord{}).
				Select("item_id, MAX(recorded_date)").
				Where("item_id IN ?", itemIDs).
				Group("item_id"),
		).
		Find(&records).Error
	return records, err
}

func (r *priceRepo) DeleteByItem(itemID uint) error {
	return r.db.Where("item_id = ?", itemID).Delete(&model.PriceRecord{}).Error
}

func (r *priceRepo) DeleteByItemAndDate(itemID uint, date string) error {
	return r.db.Where("item_id = ? AND recorded_date = ?", itemID, date).Delete(&model.PriceRecord{}).Error
}
