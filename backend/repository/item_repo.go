package repository

import (
	"time"

	"artale_market/model"

	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]model.Item, error)
	FindWithFilters(pcts []int, categories []string) ([]model.Item, error)
	FindPage(pcts []int, categories []string, itemTypes []int, sortBy string, date string, page, pageSize int) ([]model.Item, int64, error)
	FindScrollPage(pcts []int, categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error)
	FindSkillBookPage(categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error)
	FindByID(id uint) (*model.Item, error)
	FindByIDSummary(id uint, today, yesterday, threeDaysAgo string) (*model.PriceSummary, error)
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

func (r *itemRepo) FindPage(pcts []int, categories []string, itemTypes []int, sortBy string, date string, page, pageSize int) ([]model.Item, int64, error) {
	var items []model.Item
	var total int64

	applyFilters := func(q *gorm.DB) *gorm.DB {
		if len(pcts) > 0 {
			q = q.Where("percentage IN ?", pcts)
		}
		if len(categories) > 0 {
			q = q.Where("category IN ?", categories)
		}
		if len(itemTypes) > 0 {
			q = q.Where("item_type IN ?", itemTypes)
		}
		return q
	}

	if err := applyFilters(r.db.Model(&model.Item{})).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := applyFilters(r.db.Model(&model.Item{}))

	ref, _ := time.Parse("2006-01-02", date)
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")

	switch sortBy {
	case "price_desc":
		q = q.Joins("LEFT JOIN price_records pr_s ON pr_s.item_id = items.id AND pr_s.recorded_date = ?", date).
			Order("CASE WHEN pr_s.price IS NULL THEN 1 ELSE 0 END, pr_s.price DESC")
	case "price_asc":
		q = q.Joins("LEFT JOIN price_records pr_s ON pr_s.item_id = items.id AND pr_s.recorded_date = ?", date).
			Order("CASE WHEN pr_s.price IS NULL THEN 1 ELSE 0 END, pr_s.price ASC")
	case "change_desc":
		q = q.Joins("LEFT JOIN price_records pr_t ON pr_t.item_id = items.id AND pr_t.recorded_date = ?", date).
			Joins("LEFT JOIN price_records pr_y ON pr_y.item_id = items.id AND pr_y.recorded_date = ?", yesterday).
			Order("CASE WHEN pr_t.price IS NULL OR pr_y.price IS NULL OR pr_y.price = 0 THEN 1 ELSE 0 END, (pr_t.price - pr_y.price) / pr_y.price DESC")
	case "change_asc":
		q = q.Joins("LEFT JOIN price_records pr_t ON pr_t.item_id = items.id AND pr_t.recorded_date = ?", date).
			Joins("LEFT JOIN price_records pr_y ON pr_y.item_id = items.id AND pr_y.recorded_date = ?", yesterday).
			Order("CASE WHEN pr_t.price IS NULL OR pr_y.price IS NULL OR pr_y.price = 0 THEN 1 ELSE 0 END, (pr_t.price - pr_y.price) / pr_y.price ASC")
	case "percentage_asc":
		q = q.Order("items.percentage ASC, items.name ASC")
	case "percentage_desc":
		q = q.Order("items.percentage DESC, items.name ASC")
	default:
		q = q.Order("items.name ASC")
	}

	if pageSize > 0 {
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	return items, total, q.Find(&items).Error
}

func (r *itemRepo) FindScrollPage(pcts []int, categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error) {
	var total int64
	countQ := r.db.Model(&model.Item{}).Where("item_type = ?", model.ItemTypeScroll)
	if len(pcts) > 0 {
		countQ = countQ.Where("percentage IN ?", pcts)
	}
	if len(categories) > 0 {
		countQ = countQ.Where("category IN ?", categories)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := r.db.Model(&model.Item{}).
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_3days.price AS three_days_ago_price, (pr_today.price - pr_yesterday.price) / NULLIF(pr_yesterday.price, 0) AS change_pct`).
		Joins("LEFT JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?",
			today, yesterday, threeDaysAgo).
		Where("items.item_type = ?", model.ItemTypeScroll)

	if len(pcts) > 0 {
		q = q.Where("items.percentage IN ?", pcts)
	}
	if len(categories) > 0 {
		q = q.Where("items.category IN ?", categories)
	}

	switch sortBy {
	case "price_asc":
		q = q.Order("pr_today.price ASC NULLS LAST")
	case "change_desc":
		q = q.Order("change_pct DESC NULLS LAST")
	case "change_asc":
		q = q.Order("change_pct ASC NULLS LAST")
	case "percentage_asc":
		q = q.Order("items.percentage ASC, items.name ASC")
	case "percentage_desc":
		q = q.Order("items.percentage DESC, items.name ASC")
	default: // price_desc
		q = q.Order("pr_today.price DESC NULLS LAST")
	}

	if pageSize > 0 {
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	var summaries []model.PriceSummary
	return summaries, total, q.Scan(&summaries).Error
}

func (r *itemRepo) FindSkillBookPage(categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error) {
	var total int64
	countQ := r.db.Model(&model.Item{}).Where("item_type = ?", model.ItemTypeSkillBook)
	if len(categories) > 0 {
		countQ = countQ.Where("category IN ?", categories)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := r.db.Model(&model.Item{}).
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_3days.price AS three_days_ago_price, (pr_today.price - pr_yesterday.price) / NULLIF(pr_yesterday.price, 0) AS change_pct`).
		Joins("LEFT JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?",
			today, yesterday, threeDaysAgo).
		Where("items.item_type = ?", model.ItemTypeSkillBook)

	if len(categories) > 0 {
		q = q.Where("items.category IN ?", categories)
	}

	switch sortBy {
	case "price_asc":
		q = q.Order("pr_today.price ASC NULLS LAST")
	case "change_desc":
		q = q.Order("change_pct DESC NULLS LAST")
	case "change_asc":
		q = q.Order("change_pct ASC NULLS LAST")
	case "percentage_asc":
		q = q.Order("items.percentage ASC, items.name ASC")
	case "percentage_desc":
		q = q.Order("items.percentage DESC, items.name ASC")
	default: // price_desc
		q = q.Order("pr_today.price DESC NULLS LAST")
	}

	if pageSize > 0 {
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	var summaries []model.PriceSummary
	return summaries, total, q.Scan(&summaries).Error
}

func (r *itemRepo) FindByID(id uint) (*model.Item, error) {
	var item model.Item
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepo) FindByIDSummary(id uint, today, yesterday, threeDaysAgo string) (*model.PriceSummary, error) {
	var summary model.PriceSummary
	err := r.db.Model(&model.Item{}).
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_3days.price AS three_days_ago_price`).
		Joins("LEFT JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?", today, yesterday, threeDaysAgo).
		Where("items.id = ?", id).
		Scan(&summary).Error
	return &summary, err
}

func (r *itemRepo) FindTracked(date string) ([]model.Item, error) {
	var items []model.Item
	err := r.db.
		Joins("LEFT JOIN price_records ON price_records.item_id = items.id AND price_records.recorded_date = ?", date).
		Where("items.track_priority > 0 AND items.track_priority < ? AND price_records.id IS NULL", model.TrackPriorityNotSeen).
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
