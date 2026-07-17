package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindByItemType(itemType model.ItemType) ([]string, error)
	FindUsedByItemType(itemType model.ItemType) ([]string, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindByItemType(itemType model.ItemType) ([]string, error) {
	var names []string
	q := r.db.Model(&model.Category{}).Order("name")
	if itemType > 0 {
		q = q.Where("item_type = ?", itemType)
	}
	err := q.Distinct().Pluck("name", &names).Error
	return names, err
}

// FindUsedByItemType 從 items table 撈實際使用中的分類（無重複）
func (r *categoryRepo) FindUsedByItemType(itemType model.ItemType) ([]string, error) {
	var names []string
	q := r.db.Model(&model.Item{}).Where("is_hidden = false AND category != ''")
	if itemType > 0 {
		q = q.Where("item_type = ?", itemType)
	}
	err := q.Distinct().Order("category").Pluck("category", &names).Error
	return names, err
}
