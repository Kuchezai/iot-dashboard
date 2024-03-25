package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TelegramSender struct {
	chatID int
	apiKey string
}

type Message struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func NewTelegramSender(chatID int, apiKey string) *TelegramSender {
	return &TelegramSender{
		chatID: chatID,
		apiKey: apiKey,
	}
}

func (s *TelegramSender) SendMessage(messageText string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.apiKey)

	msg := &Message{
		ChatID: int(s.chatID),
		Text:   messageText,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			fmt.Println("failed to close response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send successful request. Status was %q", response.Status)
	}
	return nil
}
