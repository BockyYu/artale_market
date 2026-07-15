package service

import (
	"errors"
	"math"
	"sort"
	"time"

	"artale_market/model"
	"artale_market/repository"
)


type ItemService interface {
	GetAll() ([]model.Item, error)
	GetAllWithLatestPrice(sortBy, search string, filterType, filterPriority, page, pageSize int) (*model.AdminItemPage, error)
	GetByID(id uint) (*model.Item, error)
	GetPriceSummary(id uint) (*model.PriceSummary, error)
	GetTracked(date string) ([]model.Item, error)
	GetAllForExport(itemType int, dates [7]string) ([]model.ExportRow, error)
	GetAllForExportDynamic(itemType int, dates []string) ([]model.ExportRowDynamic, error)
	SetHidden(id uint, hidden bool) error
	Create(item *model.Item) error
	Update(id uint, name, englishName string, searchMode int, itemType model.ItemType, percentage int, category, description string) (*model.Item, error)
	SetTracked(id uint, priority model.TrackPriority) (*model.Item, error)
	Delete(id uint) error
}

type itemService struct {
	itemRepo  repository.ItemRepository
	priceRepo repository.PriceRepository
	queryRepo repository.QueryRepository
}

func NewItemService(ir repository.ItemRepository, pr repository.PriceRepository, qr repository.QueryRepository) ItemService {
	return &itemService{itemRepo: ir, priceRepo: pr, queryRepo: qr}
}

func (s *itemService) GetAll() ([]model.Item, error) {
	return s.itemRepo.FindAll()
}

func (s *itemService) GetAllWithLatestPrice(sortBy, search string, filterType, filterPriority, page, pageSize int) (*model.AdminItemPage, error) {
	viewsSort := sortBy == "views_desc" || sortBy == "views_asc"

	fetchPage, fetchSize := page, pageSize
	if viewsSort {
		fetchPage, fetchSize = 1, 0 // fetch all, paginate in memory
	}

	rows, total, err := s.itemRepo.FindAllWithLatestPrice(sortBy, search, filterType, filterPriority, fetchPage, fetchSize)
	if err != nil {
		return nil, err
	}

	today := time.Now().In(time.FixedZone("Asia/Taipei", 8*60*60)).Format("2006-01-02")
	views, _ := s.queryRepo.GetAllItemViews(today)
	for i := range rows {
		rows[i].TodayViews = views[rows[i].ID]
	}

	if viewsSort {
		sort.Slice(rows, func(i, j int) bool {
			if sortBy == "views_desc" {
				return rows[i].TodayViews > rows[j].TodayViews
			}
			return rows[i].TodayViews < rows[j].TodayViews
		})
		start := (page - 1) * pageSize
		if start > len(rows) {
			start = len(rows)
		}
		end := start + pageSize
		if end > len(rows) {
			end = len(rows)
		}
		rows = rows[start:end]
	}

	return &model.AdminItemPage{Data: rows, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *itemService) GetByID(id uint) (*model.Item, error) {
	return s.itemRepo.FindByID(id)
}

func (s *itemService) GetPriceSummary(id uint) (*model.PriceSummary, error) {
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := now.AddDate(0, 0, -3).Format("2006-01-02")

	summary, err := s.itemRepo.FindByIDSummary(id, today, yesterday, threeDaysAgo)
	if err != nil {
		return nil, err
	}
	if summary.TodayPrice != nil && summary.YesterdayPrice != nil && *summary.YesterdayPrice != 0 {
		pct := ((*summary.TodayPrice - *summary.YesterdayPrice) / *summary.YesterdayPrice) * 100
		pct = math.Round(pct*100) / 100
		summary.ChangePercent = &pct
	}
	return summary, nil
}

var ErrItemAlreadyExists = errors.New("道具名稱已存在資料庫，不可新增")

func (s *itemService) Create(item *model.Item) error {
	if existing, err := s.itemRepo.FindByName(item.Name); err == nil && existing != nil {
		return ErrItemAlreadyExists
	}
	return s.itemRepo.Create(item)
}

func (s *itemService) Update(id uint, name, englishName string, searchMode int, itemType model.ItemType, percentage int, category, description string) (*model.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err = s.itemRepo.Update(item, map[string]any{
		"name":         name,
		"english_name": englishName,
		"search_mode":  searchMode,
		"item_type":    itemType,
		"percentage":   percentage,
		"category":     category,
		"description":  description,
	}); err != nil {
		return nil, err
	}
	return s.itemRepo.FindByID(id)
}

func (s *itemService) GetTracked(date string) ([]model.Item, error) {
	return s.itemRepo.FindTracked(date)
}

func (s *itemService) GetAllForExport(itemType int, dates [7]string) ([]model.ExportRow, error) {
	return s.itemRepo.FindAllForExport(itemType, dates)
}

func (s *itemService) GetAllForExportDynamic(itemType int, dates []string) ([]model.ExportRowDynamic, error) {
	return s.itemRepo.FindAllForExportDynamic(itemType, dates)
}

func (s *itemService) SetHidden(id uint, hidden bool) error {
	return s.itemRepo.SetHidden(id, hidden)
}

func (s *itemService) SetTracked(id uint, priority model.TrackPriority) (*model.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.itemRepo.Update(item, map[string]any{"track_priority": priority}); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) Delete(id uint) error {
	if _, err := s.itemRepo.FindByID(id); err != nil {
		return err
	}
	if err := s.priceRepo.DeleteByItem(id); err != nil {
		return err
	}
	return s.itemRepo.Delete(id)
}
