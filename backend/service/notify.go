package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"artale_market/model"
)

// SendMessage dispatches a plain-text (or HTML for TG) message to the bot's platform.
func SendMessage(bot *model.NotifyBot, text string) error {
	switch bot.Platform {
	case "tg":
		return sendTGMessage(bot.Token, bot.ChatID, text)
	case "line":
		return sendLINEMessage(bot.Token, text)
	case "dc":
		return sendDiscordMessage(bot.Token, text)
	default:
		return fmt.Errorf("unsupported platform: %s", bot.Platform)
	}
}

func sendTGMessage(token, chatID, text string) error {
	if token == "" || chatID == "" {
		return fmt.Errorf("TG token 或 chat_id 未設定")
	}
	payload, _ := json.Marshal(map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	})
	resp, err := http.Post(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API 回傳 status %d：%s", resp.StatusCode, string(body))
	}
	return nil
}

func sendLINEMessage(token, text string) error {
	if token == "" {
		return fmt.Errorf("LINE token 未設定")
	}
	form := url.Values{"message": {"\n" + text}}
	req, _ := http.NewRequest(http.MethodPost,
		"https://notify-api.line.me/api/notify",
		strings.NewReader(form.Encode()),
	)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LINE Notify API 回傳 status %d", resp.StatusCode)
	}
	return nil
}

func sendDiscordMessage(webhookURL, text string) error {
	if webhookURL == "" {
		return fmt.Errorf("Discord webhook URL 未設定")
	}
	payload, _ := json.Marshal(map[string]any{"content": text})
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Discord API 回傳 status %d", resp.StatusCode)
	}
	return nil
}
