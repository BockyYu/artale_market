package service

import (
	"artale_market/repository"
)

type FrequentItem struct {
	ItemID     uint   `json:"item_id"`
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
	Category   string `json:"category"`
	Count      int    `json:"count"`
}

type QueryService interface {
	RecordQuery(userID string, itemID uint) error
	GetFrequent(userID string) ([]FrequentItem, error)
}

type queryService struct {
	queryRepo repository.QueryRepository
	itemRepo  repository.ItemRepository
}

func NewQueryService(qr repository.QueryRepository, ir repository.ItemRepository) QueryService {
	return &queryService{queryRepo: qr, itemRepo: ir}
}

func (s *queryService) RecordQuery(userID string, itemID uint) error {
	return s.queryRepo.RecordQuery(userID, itemID)
}

func (s *queryService) GetFrequent(userID string) ([]FrequentItem, error) {
	entries, err := s.queryRepo.GetFrequent(userID, 5)
	if err != nil {
		return nil, err
	}
	results := make([]FrequentItem, 0, len(entries))
	for _, e := range entries {
		item, err := s.itemRepo.FindByID(e.ItemID)
		if err != nil {
			continue // 商品已被刪除，跳過
		}
		results = append(results, FrequentItem{
			ItemID:     item.ID,
			Name:       item.Name,
			Percentage: item.Percentage,
			Category:   item.Category,
			Count:      e.Count,
		})
	}
	return results, nil
}
