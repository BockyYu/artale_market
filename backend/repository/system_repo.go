package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type SystemRepository interface {
	FindByName(name string) (*model.SystemSetting, error)
	Upsert(setting *model.SystemSetting) error
}

type systemRepository struct{ db *gorm.DB }

func NewSystemRepository(db *gorm.DB) SystemRepository {
	return &systemRepository{db}
}

func (r *systemRepository) FindByName(name string) (*model.SystemSetting, error) {
	var s model.SystemSetting
	if err := r.db.Where("name = ?", name).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *systemRepository) Upsert(setting *model.SystemSetting) error {
	return r.db.Where("name = ?", setting.Name).Assign(setting).FirstOrCreate(setting).Error
}
