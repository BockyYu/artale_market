package repository

import (
	"artale_market/model"

	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]model.Item, error)
	FindWithFilters(pcts []int, categories []string) ([]model.Item, error)
	FindByID(id uint) (*model.Item, error)
	FindTracked(date string) ([]model.Item, error)
	Create(item *model.Item) error
	Update(item *model.Item, fields map[string]any) error
	Delete(id uint) error
}

type itemRepo struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepo{db: db}
}

func (r *itemRepo) FindAll() ([]model.Item, error) {
	var items []model.Item
	err := r.db.Order("name asc").Find(&items).Error
	return items, err
}

func (r *itemRepo) FindWithFilters(pcts []int, categories []string) ([]model.Item, error) {
	var items []model.Item
	q := r.db.Order("name asc")
	if len(pcts) > 0 {
		q = q.Where("percentage IN ?", pcts)
	}
	if len(categories) > 0 {
		q = q.Where("category IN ?", categories)
	}
	return items, q.Find(&items).Error
}

func (r *itemRepo) FindByID(id uint) (*model.Item, error) {
	var item model.Item
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) FindTracked(date string) ([]model.Item, error) {
	var items []model.Item
	err := r.db.
		Joins("LEFT JOIN price_records ON price_records.item_id = items.id AND price_records.recorded_date = ?", date).
		Where("items.track_priority > 0 AND price_records.id IS NULL").
		Order("items.track_priority asc, items.name asc").
		Find(&items).Error
	return items, err
}

func (r *itemRepo) Create(item *model.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepo) Update(item *model.Item, fields map[string]any) error {
	return r.db.Model(item).Updates(fields).Error
}

func (r *itemRepo) Delete(id uint) error {
	return r.db.Delete(&model.Item{}, id).Error
}
