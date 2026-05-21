package service

import (
	"artale_market/model"
	"artale_market/repository"
)

type ItemService interface {
	GetAll() ([]model.Item, error)
	GetByID(id uint) (*model.Item, error)
	GetTracked(date string) ([]model.Item, error)
	Create(item *model.Item) error
	Update(id uint, name string, percentage int, category, description string) (*model.Item, error)
	SetTracked(id uint, priority int) (*model.Item, error)
	Delete(id uint) error
}

type itemService struct {
	itemRepo  repository.ItemRepository
	priceRepo repository.PriceRepository
}

func NewItemService(ir repository.ItemRepository, pr repository.PriceRepository) ItemService {
	return &itemService{itemRepo: ir, priceRepo: pr}
}

func (s *itemService) GetAll() ([]model.Item, error) {
	return s.itemRepo.FindAll()
}

func (s *itemService) GetByID(id uint) (*model.Item, error) {
	return s.itemRepo.FindByID(id)
}

func (s *itemService) Create(item *model.Item) error {
	return s.itemRepo.Create(item)
}

func (s *itemService) Update(id uint, name string, percentage int, category, description string) (*model.Item, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	err = s.itemRepo.Update(item, map[string]any{
		"name":        name,
		"percentage":  percentage,
		"category":    category,
		"description": description,
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) GetTracked(date string) ([]model.Item, error) {
	return s.itemRepo.FindTracked(date)
}

func (s *itemService) SetTracked(id uint, priority int) (*model.Item, error) {
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
