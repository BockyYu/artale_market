package repository

import (
	"artale_market/model"
	"time"

	"gorm.io/gorm"
)

type AlertRepository interface {
	List() ([]model.PriceAlert, error)
	Create(alert *model.PriceAlert) error
	Delete(id uint) error
	ToggleActive(id uint, isActive bool) error
	FindActiveByItem(itemID uint) ([]model.PriceAlert, error)
	FindAllActive() ([]model.PriceAlert, error)
	UpdateLastTriggered(id uint) error
}

type alertRepository struct {
	db *gorm.DB
}

func NewAlertRepository(db *gorm.DB) AlertRepository {
	return &alertRepository{db: db}
}

func (r *alertRepository) List() ([]model.PriceAlert, error) {
	var alerts []model.PriceAlert
	err := r.db.Preload("Item").Preload("Bot").Order("created_at desc").Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) Create(alert *model.PriceAlert) error {
	return r.db.Create(alert).Error
}

func (r *alertRepository) Delete(id uint) error {
	return r.db.Delete(&model.PriceAlert{}, id).Error
}

func (r *alertRepository) ToggleActive(id uint, isActive bool) error {
	return r.db.Model(&model.PriceAlert{}).Where("id = ?", id).Update("is_active", isActive).Error
}

func (r *alertRepository) FindActiveByItem(itemID uint) ([]model.PriceAlert, error) {
	var alerts []model.PriceAlert
	err := r.db.Preload("Bot").Where("item_id = ? AND is_active = true", itemID).Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) FindAllActive() ([]model.PriceAlert, error) {
	var alerts []model.PriceAlert
	err := r.db.Preload("Item").Preload("Bot").Where("is_active = true").Find(&alerts).Error
	return alerts, err
}

func (r *alertRepository) UpdateLastTriggered(id uint) error {
	now := time.Now()
	return r.db.Model(&model.PriceAlert{}).Where("id = ?", id).Update("last_triggered_at", now).Error
}
