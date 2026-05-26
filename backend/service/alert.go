package service

import (
	"artale_market/model"
	"artale_market/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type AlertService interface {
	List() ([]model.PriceAlert, error)
	Create(itemID uint, botID *uint, threshold float64, note string) (*model.PriceAlert, error)
	Delete(id uint) error
	ToggleActive(id uint, isActive bool) error
	CheckAndNotify(itemID uint, itemName string, price float64)
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

func (s *alertService) Delete(id uint) error {
	return s.alertRepo.Delete(id)
}

func (s *alertService) ToggleActive(id uint, isActive bool) error {
	return s.alertRepo.ToggleActive(id, isActive)
}

func (s *alertService) CheckAndNotify(itemID uint, itemName string, price float64) {
	alerts, err := s.alertRepo.FindActiveByItem(itemID)
	if err != nil {
		log.Printf("[Alert] failed to fetch alerts for item %d: %v", itemID, err)
		return
	}
	for _, alert := range alerts {
		if price > alert.ThresholdPrice {
			continue
		}
		// 同一提醒每小時最多通知一次，避免洗版
		if alert.LastTriggeredAt != nil && time.Since(*alert.LastTriggeredAt) < time.Hour {
			continue
		}
		go func(a model.PriceAlert) {
			sendNotification(a.Bot, itemName, price, a.ThresholdPrice, a.Note)
			_ = s.alertRepo.UpdateLastTriggered(a.ID)
		}(alert)
	}
}

// sendNotification 根據機器人平台發送通知，無機器人時回退到環境變數 TG 設定
func sendNotification(bot *model.NotifyBot, itemName string, price, threshold float64, note string) {
	if bot == nil {
		sendTG(os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"), itemName, price, threshold, note)
		return
	}
	switch bot.Platform {
	case "tg":
		sendTG(bot.Token, bot.ChatID, itemName, price, threshold, note)
	case "line":
		sendLINE(bot.Token, itemName, price, threshold, note)
	case "dc":
		sendDiscord(bot.Token, itemName, price, threshold, note)
	}
}

func buildText(itemName string, price, threshold float64, note string) string {
	text := fmt.Sprintf("🔔 價格提醒\n📦 %s\n💰 目前價格：%s\n🎯 觸發門檻：%s",
		itemName, fmtPrice(price), fmtPrice(threshold))
	if note != "" {
		text += "\n📝 備註：" + note
	}
	return text
}

func sendTG(token, chatID, itemName string, price, threshold float64, note string) {
	if token == "" || chatID == "" {
		return
	}
	htmlText := fmt.Sprintf("🔔 <b>價格提醒</b>\n📦 <b>%s</b>\n💰 目前價格：<b>%s</b>\n🎯 觸發門檻：%s",
		itemName, fmtPrice(price), fmtPrice(threshold))
	if note != "" {
		htmlText += "\n📝 備註：" + note
	}
	payload, _ := json.Marshal(map[string]any{
		"chat_id":    chatID,
		"text":       htmlText,
		"parse_mode": "HTML",
	})
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		log.Printf("[Alert][TG] send failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Alert][TG] status %d", resp.StatusCode)
	}
}

func sendLINE(token, itemName string, price, threshold float64, note string) {
	if token == "" {
		return
	}
	text := buildText(itemName, price, threshold, note)
	form := url.Values{"message": {"\n" + text}}
	req, _ := http.NewRequest(http.MethodPost,
		"https://notify-api.line.me/api/notify",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[Alert][LINE] send failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Alert][LINE] status %d", resp.StatusCode)
	}
}

func sendDiscord(webhookURL, itemName string, price, threshold float64, note string) {
	if webhookURL == "" {
		return
	}
	text := buildText(itemName, price, threshold, note)
	payload, _ := json.Marshal(map[string]any{"content": text})
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Printf("[Alert][DC] send failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("[Alert][DC] status %d", resp.StatusCode)
	}
}

func fmtPrice(p float64) string {
	s := fmt.Sprintf("%.0f", p)
	n := len(s)
	if n <= 3 {
		return s
	}
	result := make([]byte, 0, n+(n-1)/3)
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
