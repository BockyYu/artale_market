package service

import (
	"math"
	"time"

	"artale_market/model"
	"artale_market/repository"
)

type PriceService interface {
	GetSummary(date string, pcts []int, categories []string) ([]model.PriceSummary, error)
	Record(itemID uint, price float64) (*model.PriceRecord, error)
	GetHistory(itemID uint) ([]model.PriceRecord, error)
}

type priceService struct {
	itemRepo  repository.ItemRepository
	priceRepo repository.PriceRepository
}

func NewPriceService(ir repository.ItemRepository, pr repository.PriceRepository) PriceService {
	return &priceService{itemRepo: ir, priceRepo: pr}
}

func (svc *priceService) GetSummary(date string, pcts []int, categories []string) ([]model.PriceSummary, error) {
	var items []model.Item
	var err error
	if len(pcts) > 0 || len(categories) > 0 {
		items, err = svc.itemRepo.FindWithFilters(pcts, categories)
	} else {
		items, err = svc.itemRepo.FindAll()
	}
	if err != nil {
		return nil, err
	}

	ref, err := time.Parse("2006-01-02", date)
	if err != nil {
		ref = time.Now()
	}
	today := ref.Format("2006-01-02")
	weekAgo := ref.AddDate(0, 0, -7).Format("2006-01-02")

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

		if r, err := svc.priceRepo.FindByItemAndDate(item.ID, today); err == nil {
			p := r.Price
			s.TodayPrice = &p
		}

		if w, err := svc.priceRepo.FindByItemAndDate(item.ID, weekAgo); err == nil {
			p := w.Price
			s.WeekAgoPrice = &p
		}

		if s.TodayPrice != nil && s.WeekAgoPrice != nil && *s.WeekAgoPrice != 0 {
			pct := ((*s.TodayPrice - *s.WeekAgoPrice) / *s.WeekAgoPrice) * 100
			pct = math.Round(pct*100) / 100
			s.ChangePercent = &pct
		}

		summaries = append(summaries, s)
	}

	return summaries, nil
}

func (svc *priceService) Record(itemID uint, price float64) (*model.PriceRecord, error) {
	if _, err := svc.itemRepo.FindByID(itemID); err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")

	existing, err := svc.priceRepo.FindByItemAndDate(itemID, today)
	if err != nil {
		record := &model.PriceRecord{
			ItemID:       itemID,
			Price:        price,
			RecordedDate: time.Now(),
		}
		if err := svc.priceRepo.Create(record); err != nil {
			return nil, err
		}
		return record, nil
	}

	if err := svc.priceRepo.UpdatePrice(existing, price); err != nil {
		return nil, err
	}
	existing.Price = price
	return existing, nil
}

func (svc *priceService) GetHistory(itemID uint) ([]model.PriceRecord, error) {
	return svc.priceRepo.FindHistoryByItem(itemID, 30)
}
