package service

import (
	"fmt"
	"math"
	"time"

	"artale_market/model"
	"artale_market/repository"
)

type PriceService interface {
	GetSummary(date string, pcts []int, categories []string, itemTypes []int, sortBy string, page, pageSize int) (*model.PagedSummary, error)
	GetScrollSummary(date string, pcts []int, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error)
	GetSkillBookSummary(date string, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error)
	GetEquipSummary(date string, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error)
	GetOtherSummary(date string, sortBy string, page, pageSize int) (*model.PagedSummary, error)
	Record(itemID uint, price float64, date string, source string) (*model.PriceRecord, error)
	GetLatest(itemID uint) (*model.PriceRecord, error)
	GetLatestBatch(itemIDs []uint) ([]model.PriceRecord, error)
	GetHistory(itemID uint) ([]model.PriceRecord, error)
	GetAllHistory(itemID uint) ([]model.PriceRecord, error)
	GetPriceHistories(itemID uint) ([]model.PriceHistory, error)
	DeletePriceHistory(id uint) error
	TogglePriceHistoryHidden(id uint, isHidden bool) error
}

type priceService struct {
	itemRepo    repository.ItemRepository
	priceRepo   repository.PriceRepository
	historyRepo repository.PriceHistoryRepository
}

func NewPriceService(ir repository.ItemRepository, pr repository.PriceRepository, hr repository.PriceHistoryRepository) PriceService {
	return &priceService{itemRepo: ir, priceRepo: pr, historyRepo: hr}
}

func (svc *priceService) GetSummary(date string, pcts []int, categories []string, itemTypes []int, sortBy string, page, pageSize int) (*model.PagedSummary, error) {
	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := ref.AddDate(0, 0, -3).Format("2006-01-02")

	items, total, err := svc.itemRepo.FindPage(pcts, categories, itemTypes, sortBy, today, page, pageSize)
	if err != nil {
		return nil, err
	}

	itemIDs := make([]uint, len(items))
	for i, item := range items {
		itemIDs[i] = item.ID
	}

	priceRecords, err := svc.priceRepo.FindByItemsAndDates(itemIDs, []string{today, yesterday, threeDaysAgo})
	if err != nil {
		return nil, err
	}

	// key: "itemID|date"
	priceMap := make(map[string]*model.PriceRecord, len(priceRecords))
	for i := range priceRecords {
		r := &priceRecords[i]
		key := fmt.Sprintf("%d|%s", r.ItemID, r.RecordedDate.Format("2006-01-02"))
		priceMap[key] = r
	}

	summaries := make([]model.PriceSummary, 0, len(items))
	for _, item := range items {
		s := model.PriceSummary{
			ItemID:         item.ID,
			ItemName:       item.Name,
			ItemPercentage: item.Percentage,
			ItemType:       item.ItemType,
			Category:       item.Category,
			Description:    item.Description,
		}

		if r := priceMap[fmt.Sprintf("%d|%s", item.ID, today)]; r != nil {
			p := r.Price
			s.TodayPrice = &p
			s.TodayCreatedAt = &r.CreatedAt
			if !r.UpdatedAt.IsZero() {
				s.TodayUpdatedAt = &r.UpdatedAt
			}
		}
		if r := priceMap[fmt.Sprintf("%d|%s", item.ID, yesterday)]; r != nil {
			p := r.Price
			s.YesterdayPrice = &p
		}
		if r := priceMap[fmt.Sprintf("%d|%s", item.ID, threeDaysAgo)]; r != nil {
			p := r.Price
			s.ThreeDaysAgoPrice = &p
		}

		if s.TodayPrice != nil && s.YesterdayPrice != nil && *s.YesterdayPrice != 0 {
			pct := ((*s.TodayPrice - *s.YesterdayPrice) / *s.YesterdayPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}

		summaries = append(summaries, s)
	}

	return &model.PagedSummary{
		Data:     summaries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (svc *priceService) GetScrollSummary(date string, pcts []int, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error) {
	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := ref.AddDate(0, 0, -3).Format("2006-01-02")

	summaries, total, err := svc.itemRepo.FindScrollPage(pcts, categories, sortBy, today, yesterday, threeDaysAgo, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range summaries {
		s := &summaries[i]
		if s.TodayPrice != nil && s.YesterdayPrice != nil && *s.YesterdayPrice != 0 {
			pct := ((*s.TodayPrice - *s.YesterdayPrice) / *s.YesterdayPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}
	}

	return &model.PagedSummary{Data: summaries, Total: total, Page: page, PageSize: pageSize}, nil
}

func (svc *priceService) GetSkillBookSummary(date string, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error) {
	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := ref.AddDate(0, 0, -3).Format("2006-01-02")

	summaries, total, err := svc.itemRepo.FindSkillBookPage(categories, sortBy, today, yesterday, threeDaysAgo, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range summaries {
		s := &summaries[i]
		if s.TodayPrice != nil && s.YesterdayPrice != nil && *s.YesterdayPrice != 0 {
			pct := ((*s.TodayPrice - *s.YesterdayPrice) / *s.YesterdayPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}
	}

	return &model.PagedSummary{
		Data:     summaries,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (svc *priceService) GetEquipSummary(date string, categories []string, sortBy string, page, pageSize int) (*model.PagedSummary, error) {
	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := ref.AddDate(0, 0, -3).Format("2006-01-02")

	summaries, total, err := svc.itemRepo.FindEquipPage(categories, sortBy, today, yesterday, threeDaysAgo, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range summaries {
		s := &summaries[i]
		if s.TodayPrice != nil && s.YesterdayPrice != nil && *s.YesterdayPrice != 0 {
			pct := ((*s.TodayPrice - *s.YesterdayPrice) / *s.YesterdayPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}
	}

	return &model.PagedSummary{Data: summaries, Total: total, Page: page, PageSize: pageSize}, nil
}

func (svc *priceService) GetOtherSummary(date string, sortBy string, page, pageSize int) (*model.PagedSummary, error) {
	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	yesterday := ref.AddDate(0, 0, -1).Format("2006-01-02")
	threeDaysAgo := ref.AddDate(0, 0, -3).Format("2006-01-02")

	summaries, total, err := svc.itemRepo.FindOtherPage(sortBy, today, yesterday, threeDaysAgo, page, pageSize)
	if err != nil {
		return nil, err
	}

	for i := range summaries {
		s := &summaries[i]
		if s.TodayPrice != nil && s.YesterdayPrice != nil && *s.YesterdayPrice != 0 {
			pct := ((*s.TodayPrice - *s.YesterdayPrice) / *s.YesterdayPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}
	}

	return &model.PagedSummary{Data: summaries, Total: total, Page: page, PageSize: pageSize}, nil
}

func (svc *priceService) Record(itemID uint, price float64, date string, source string) (*model.PriceRecord, error) {
	if _, err := svc.itemRepo.FindByID(itemID); err != nil {
		return nil, err
	}

	recordedDate, _ := time.Parse("2006-01-02", date)

	var record *model.PriceRecord
	existing, err := svc.priceRepo.FindByItemAndDate(itemID, date)
	if err != nil {
		record = &model.PriceRecord{
			ItemID:       itemID,
			Price:        price,
			RecordedDate: recordedDate,
		}
		if err := svc.priceRepo.Create(record); err != nil {
			return nil, err
		}
	} else {
		if existing.Price == price {
			return existing, nil // 價格與今日紀錄相同，略過更新
		}
		if err := svc.priceRepo.UpdatePrice(existing, price); err != nil {
			return nil, err
		}
		existing.Price = price
		existing.UpdatedAt = time.Now()
		record = existing
	}

	_ = svc.historyRepo.Create(&model.PriceHistory{ItemID: itemID, Price: price, Source: source})

	return record, nil
}

func (svc *priceService) GetPriceHistories(itemID uint) ([]model.PriceHistory, error) {
	return svc.historyRepo.FindByItem(itemID)
}

func (svc *priceService) DeletePriceHistory(id uint) error {
	h, err := svc.historyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if err := svc.historyRepo.Delete(id); err != nil {
		return err
	}
	return svc.syncPriceRecord(h.ItemID, h.RecordedAt)
}

func (svc *priceService) TogglePriceHistoryHidden(id uint, isHidden bool) error {
	h, err := svc.historyRepo.FindByID(id)
	if err != nil {
		return err
	}
	if err := svc.historyRepo.ToggleHidden(id, isHidden); err != nil {
		return err
	}
	return svc.syncPriceRecord(h.ItemID, h.RecordedAt)
}

// syncPriceRecord 在 history 異動後，重新計算並同步當日 price_record。
// 若當日已無可見 history，直接刪除 price_record（前後台不再顯示）。
func (svc *priceService) syncPriceRecord(itemID uint, recordedAt time.Time) error {
	loc := time.FixedZone("Asia/Taipei", 8*60*60)
	date := recordedAt.In(loc).Format("2006-01-02")

	remaining, err := svc.historyRepo.FindVisibleByItemAndDate(itemID, date)
	if err != nil {
		return err
	}

	if len(remaining) == 0 {
		return svc.priceRepo.DeleteByItemAndDate(itemID, date)
	}

	// 取剩餘可見 history 的最低價
	minPrice := remaining[0].Price
	for _, h := range remaining[1:] {
		if h.Price < minPrice {
			minPrice = h.Price
		}
	}

	record, err := svc.priceRepo.FindByItemAndDate(itemID, date)
	if err != nil {
		return nil // price_record 不存在，無需同步
	}
	return svc.priceRepo.UpdatePrice(record, minPrice)
}

func (svc *priceService) GetLatest(itemID uint) (*model.PriceRecord, error) {
	return svc.priceRepo.FindLatestByItem(itemID)
}

func (svc *priceService) GetLatestBatch(itemIDs []uint) ([]model.PriceRecord, error) {
	return svc.priceRepo.FindLatestByItems(itemIDs)
}

func (svc *priceService) GetHistory(itemID uint) ([]model.PriceRecord, error) {
	return svc.priceRepo.FindHistoryByItem(itemID, 30)
}

func (svc *priceService) GetAllHistory(itemID uint) ([]model.PriceRecord, error) {
	return svc.priceRepo.FindHistoryByItem(itemID, -1)
}
