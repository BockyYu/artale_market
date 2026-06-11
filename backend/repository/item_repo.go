package repository

import (
	"time"

	"artale_market/model"

	"gorm.io/gorm"
)

type ItemRepository interface {
	FindAll() ([]model.Item, error)
	FindAllWithLatestPrice(sortBy, search string, filterType, filterPriority, page, pageSize int) ([]model.ItemAdminRow, int64, error)
	FindWithFilters(pcts []int, categories []string) ([]model.Item, error)
	FindPage(pcts []int, categories []string, itemTypes []int, sortBy string, date string, page, pageSize int) ([]model.Item, int64, error)
	FindScrollPage(pcts []int, categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error)
	FindSkillBookPage(categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error)
	FindEquipPage(categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error)
	FindByID(id uint) (*model.Item, error)
	FindByName(name string) (*model.Item, error)
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

func (r *itemRepo) applyAdminFilters(q *gorm.DB, search string, filterType, filterPriority int) *gorm.DB {
	if search != "" {
		q = q.Where("items.name ILIKE ? OR items.category ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if filterType > 0 {
		q = q.Where("items.item_type = ?", filterType)
	}
	if filterPriority >= 0 {
		q = q.Where("items.track_priority = ?", filterPriority)
	}
	return q
}

func (r *itemRepo) FindAllWithLatestPrice(sortBy, search string, filterType, filterPriority, page, pageSize int) ([]model.ItemAdminRow, int64, error) {
	var rows []model.ItemAdminRow
	var total int64

	order := "items.id ASC"
	switch sortBy {
	case "id_desc":
		order = "items.id DESC"
	case "price_desc":
		order = "latest_price DESC NULLS LAST, items.id ASC"
	case "price_asc":
		order = "latest_price ASC NULLS LAST, items.id ASC"
	case "changes_desc":
		order = "today_changes DESC NULLS LAST, items.id ASC"
	case "changes_asc":
		order = "today_changes ASC NULLS LAST, items.id ASC"
	}

	if err := r.applyAdminFilters(r.db.Model(&model.Item{}), search, filterType, filterPriority).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	loc, _ := time.LoadLocation("Asia/Taipei")
	today := time.Now().In(loc).Format("2006-01-02")
	q := r.applyAdminFilters(r.db.Model(&model.Item{}), search, filterType, filterPriority).
		Select("items.*, " +
			"(SELECT price FROM price_records WHERE item_id = items.id ORDER BY recorded_date DESC, updated_at DESC LIMIT 1) AS latest_price, " +
			"(SELECT COALESCE(NULLIF(updated_at, '0001-01-01'), created_at) FROM price_records WHERE item_id = items.id ORDER BY recorded_date DESC, updated_at DESC LIMIT 1) AS latest_price_at, " +
			"(SELECT COUNT(*) FROM price_histories WHERE item_id = items.id AND source = 'admin' AND (recorded_at AT TIME ZONE 'Asia/Taipei')::date = '" + today + "') AS today_changes").
		Order(order)

	if pageSize > 0 {
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	return rows, total, q.Scan(&rows).Error
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

	hasToday := r.db.Model(&model.Item{}).
		Joins("JOIN price_records pr_exists ON pr_exists.item_id = items.id AND pr_exists.recorded_date = ?", date)
	if err := applyFilters(hasToday).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := applyFilters(r.db.Model(&model.Item{})).
		Joins("JOIN price_records pr_exists ON pr_exists.item_id = items.id AND pr_exists.recorded_date = ?", date)

	ref, _ := time.Parse("2006-01-02", date)
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")

	switch sortBy {
	case "price_desc":
		q = q.Order("pr_exists.price DESC")
	case "price_asc":
		q = q.Order("pr_exists.price ASC")
	case "change_desc":
		q = q.Joins("LEFT JOIN price_records pr_y ON pr_y.item_id = items.id AND pr_y.recorded_date = ?", yesterday).
			Order("CASE WHEN pr_y.price IS NULL OR pr_y.price = 0 THEN 1 ELSE 0 END, (pr_exists.price - pr_y.price) / pr_y.price DESC")
	case "change_asc":
		q = q.Joins("LEFT JOIN price_records pr_y ON pr_y.item_id = items.id AND pr_y.recorded_date = ?", yesterday).
			Order("CASE WHEN pr_y.price IS NULL OR pr_y.price = 0 THEN 1 ELSE 0 END, (pr_exists.price - pr_y.price) / pr_y.price ASC")
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
	countQ := r.db.Model(&model.Item{}).Where("item_type = ?", model.ItemTypeScroll).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ?", today)
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
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_yesterday.created_at AS yesterday_created_at, pr_yesterday.updated_at AS yesterday_updated_at, pr_3days.price AS three_days_ago_price, (pr_today.price - pr_yesterday.price) / NULLIF(pr_yesterday.price, 0) AS change_pct`).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?",
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
	case "yesterday_price_desc":
		q = q.Order("pr_yesterday.price DESC NULLS LAST")
	case "yesterday_price_asc":
		q = q.Order("pr_yesterday.price ASC NULLS LAST")
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
	countQ := r.db.Model(&model.Item{}).Where("item_type = ?", model.ItemTypeSkillBook).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ?", today)
	if len(categories) > 0 {
		countQ = countQ.Where("category IN ?", categories)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := r.db.Model(&model.Item{}).
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_yesterday.created_at AS yesterday_created_at, pr_yesterday.updated_at AS yesterday_updated_at, pr_3days.price AS three_days_ago_price, (pr_today.price - pr_yesterday.price) / NULLIF(pr_yesterday.price, 0) AS change_pct`).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?",
			today, yesterday, threeDaysAgo).
		Where("items.item_type = ?", model.ItemTypeSkillBook)

	if len(categories) > 0 {
		q = q.Where("items.category IN ?", categories)
	}

	switch sortBy {
	case "price_asc":
		q = q.Order("pr_today.price ASC NULLS LAST")
	case "yesterday_price_desc":
		q = q.Order("pr_yesterday.price DESC NULLS LAST")
	case "yesterday_price_asc":
		q = q.Order("pr_yesterday.price ASC NULLS LAST")
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

func (r *itemRepo) FindEquipPage(categories []string, sortBy string, today, yesterday, threeDaysAgo string, page, pageSize int) ([]model.PriceSummary, int64, error) {
	var total int64
	countQ := r.db.Model(&model.Item{}).Where("item_type = ?", model.ItemTypeEquip).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ?", today)
	if len(categories) > 0 {
		countQ = countQ.Where("category IN ?", categories)
	}
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := r.db.Model(&model.Item{}).
		Select(`items.id AS item_id, items.name AS item_name, items.percentage AS item_percentage, items.item_type AS item_type, items.category AS category, items.description AS description, pr_today.price AS today_price, pr_today.created_at AS today_created_at, pr_today.updated_at AS today_updated_at, pr_yesterday.price AS yesterday_price, pr_yesterday.created_at AS yesterday_created_at, pr_yesterday.updated_at AS yesterday_updated_at, pr_3days.price AS three_days_ago_price, (pr_today.price - pr_yesterday.price) / NULLIF(pr_yesterday.price, 0) AS change_pct`).
		Joins("JOIN price_records pr_today ON pr_today.item_id = items.id AND pr_today.recorded_date = ? LEFT JOIN price_records pr_yesterday ON pr_yesterday.item_id = items.id AND pr_yesterday.recorded_date = ? LEFT JOIN price_records pr_3days ON pr_3days.item_id = items.id AND pr_3days.recorded_date = ?",
			today, yesterday, threeDaysAgo).
		Where("items.item_type = ?", model.ItemTypeEquip)

	if len(categories) > 0 {
		q = q.Where("items.category IN ?", categories)
	}

	switch sortBy {
	case "price_asc":
		q = q.Order("pr_today.price ASC NULLS LAST")
	case "yesterday_price_desc":
		q = q.Order("pr_yesterday.price DESC NULLS LAST")
	case "yesterday_price_asc":
		q = q.Order("pr_yesterday.price ASC NULLS LAST")
	case "change_desc":
		q = q.Order("change_pct DESC NULLS LAST")
	case "change_asc":
		q = q.Order("change_pct ASC NULLS LAST")
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

func (r *itemRepo) FindByName(name string) (*model.Item, error) {
	var item model.Item
	if err := r.db.Where("name ILIKE ?", name).First(&item).Error; err != nil {
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
