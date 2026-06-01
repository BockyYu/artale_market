package service

import (
	"artale_market/model"
	"artale_market/repository"
)

type AlertService interface {
	List() ([]model.PriceAlert, error)
	Create(itemID uint, botID *uint, threshold float64, note string) (*model.PriceAlert, error)
	Update(id uint, botID *uint, threshold float64, note string) error
	Delete(id uint) error
	ToggleActive(id uint, isActive bool) error
}

type alertService struct {
	alertRepo repository.AlertRepository
}

func NewAlertService(alertRepo repository.AlertRepository) AlertService {
	return &alertService{alertRepo: alertRepo}
}

func (s *alertService) List() ([]model.PriceAlert, error) {
	return s.alertRepo.List()
}

func (s *alertService) Create(itemID uint, botID *uint, threshold float64, note string) (*model.PriceAlert, error) {
	alert := &model.PriceAlert{
		ItemID:         itemID,
		BotID:          botID,
		ThresholdPrice: threshold,
		Note:           note,
		IsActive:       true,
	}
	if err := s.alertRepo.Create(alert); err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *alertService) Update(id uint, botID *uint, threshold float64, note string) error {
	fields := map[string]any{
		"threshold_price": threshold,
		"bot_id":          botID,
		"note":            note,
	}
	return s.alertRepo.Update(id, fields)
}

func (s *alertService) Delete(id uint) error {
	return s.alertRepo.Delete(id)
}

func (s *alertService) ToggleActive(id uint, isActive bool) error {
	return s.alertRepo.ToggleActive(id, isActive)
}

