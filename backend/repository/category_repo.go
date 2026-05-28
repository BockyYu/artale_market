package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByItemType(itemType int) ([]string, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindByItemType(itemType int) ([]string, error) {
	var names []string
	q := r.db.Model(&model.Category{}).Order("name")
	if itemType > 0 {
		q = q.Where("item_type = ?", itemType)
	}
	err := q.Pluck("name", &names).Error
	return names, err
}
