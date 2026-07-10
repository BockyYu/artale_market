package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type BotRepository interface {
	List() ([]model.NotifyBot, error)
	ListActive() ([]model.NotifyBot, error)
	FindByID(id uint) (*model.NotifyBot, error)
	Create(bot *model.NotifyBot) error
	Update(id uint, fields map[string]any) error
	Delete(id uint) error
	ToggleActive(id uint, isActive bool) error
}

type botRepository struct {
	db *gorm.DB
}

func NewBotRepository(db *gorm.DB) BotRepository {
	return &botRepository{db: db}
}

func (r *botRepository) List() ([]model.NotifyBot, error) {
	var bots []model.NotifyBot
	err := r.db.Order("created_at desc").Find(&bots).Error
	return bots, err
}

func (r *botRepository) ListActive() ([]model.NotifyBot, error) {
	var bots []model.NotifyBot
	err := r.db.Where("is_active = true").Order("created_at desc").Find(&bots).Error
	return bots, err
}

func (r *botRepository) FindByID(id uint) (*model.NotifyBot, error) {
	var bot model.NotifyBot
	err := r.db.First(&bot, id).Error
	return &bot, err
}

func (r *botRepository) Create(bot *model.NotifyBot) error {
	return r.db.Create(bot).Error
}

func (r *botRepository) Update(id uint, fields map[string]any) error {
	return r.db.Model(&model.NotifyBot{}).Where("id = ?", id).Updates(fields).Error
}

func (r *botRepository) Delete(id uint) error {
	return r.db.Delete(&model.NotifyBot{}, id).Error
}

func (r *botRepository) ToggleActive(id uint, isActive bool) error {
	return r.db.Model(&model.NotifyBot{}).Where("id = ?", id).Update("is_active", isActive).Error
}
