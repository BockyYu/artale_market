package service

import (
	"artale_market/model"
	"artale_market/repository"
)

type BotService interface {
	List() ([]model.NotifyBot, error)
	Create(name, platform, token, chatID string) (*model.NotifyBot, error)
	Update(id uint, name, platform, token, chatID string) error
	Delete(id uint) error
	ToggleActive(id uint, isActive bool) error
	SendMessage(id uint, message string) error
}

type botService struct {
	repo repository.BotRepository
}

func NewBotService(repo repository.BotRepository) BotService {
	return &botService{repo: repo}
}

func (s *botService) List() ([]model.NotifyBot, error) {
	return s.repo.List()
}

func (s *botService) Create(name, platform, token, chatID string) (*model.NotifyBot, error) {
	bot := &model.NotifyBot{
		Name:     name,
		Platform: platform,
		Token:    token,
		ChatID:   chatID,
		IsActive: true,
	}
	if err := s.repo.Create(bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func (s *botService) Update(id uint, name, platform, token, chatID string) error {
	fields := map[string]any{
		"name":     name,
		"platform": platform,
		"chat_id":  chatID,
	}
	if token != "" {
		fields["token"] = token
	}
	return s.repo.Update(id, fields)
}

func (s *botService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *botService) ToggleActive(id uint, isActive bool) error {
	return s.repo.ToggleActive(id, isActive)
}

func (s *botService) SendMessage(id uint, message string) error {
	bot, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return SendMessage(bot, message)
}
