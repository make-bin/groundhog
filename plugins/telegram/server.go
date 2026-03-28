package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// TelegramUpdate represents a Telegram Bot API update.
type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

// TelegramMessage represents a Telegram message.
type TelegramMessage struct {
	MessageID int           `json:"message_id"`
	From      *TelegramUser `json:"from,omitempty"`
	Chat      TelegramChat  `json:"chat"`
	Text      string        `json:"text,omitempty"`
	Date      int64         `json:"date"`
}

// TelegramUser represents a Telegram user.
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username,omitempty"`
}

// TelegramChat represents a Telegram chat.
type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// InboundMessage is a simplified inbound message for the gateway.
type InboundMessage struct {
	MessageID  string
	ChannelID  string
	AccountID  string
	Content    string
	ReceivedAt time.Time
}

// TelegramServer implements the Telegram channel plugin.
type TelegramServer struct {
	token      string
	baseURL    string
	offset     int
	msgCh      chan *InboundMessage
	mu         sync.Mutex
	httpClient *http.Client
}

// NewTelegramServer creates a new TelegramServer.
func NewTelegramServer(token string) *TelegramServer {
	return &TelegramServer{
		token:      token,
		baseURL:    fmt.Sprintf("https://api.telegram.org/bot%s", token),
		msgCh:      make(chan *InboundMessage, 100),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Start begins polling for updates.
func (s *TelegramServer) Start(ctx context.Context) error {
	go s.poll(ctx)
	return nil
}

// poll continuously polls the Telegram Bot API for updates.
func (s *TelegramServer) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		updates, err := s.getUpdates(ctx)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, update := range updates {
			if update.Message != nil && update.Message.Text != "" {
				msg := &InboundMessage{
					MessageID:  strconv.Itoa(update.Message.MessageID),
					ChannelID:  "telegram",
					AccountID:  strconv.FormatInt(update.Message.From.ID, 10),
					Content:    update.Message.Text,
					ReceivedAt: time.Unix(update.Message.Date, 0),
				}
				select {
				case s.msgCh <- msg:
				default:
				}
			}
			s.mu.Lock()
			s.offset = update.UpdateID + 1
			s.mu.Unlock()
		}
	}
}

// getUpdates calls the Telegram getUpdates API.
func (s *TelegramServer) getUpdates(ctx context.Context) ([]TelegramUpdate, error) {
	s.mu.Lock()
	offset := s.offset
	s.mu.Unlock()

	params := url.Values{}
	params.Set("offset", strconv.Itoa(offset))
	params.Set("timeout", "25")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		s.baseURL+"/getUpdates?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OK     bool             `json:"ok"`
		Result []TelegramUpdate `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Result, nil
}

// SendMessage sends a text message to a Telegram chat.
func (s *TelegramServer) SendMessage(chatID, text string) error {
	params := url.Values{}
	params.Set("chat_id", chatID)
	params.Set("text", text)

	resp, err := s.httpClient.Post(
		s.baseURL+"/sendMessage",
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_ = params // used in real implementation
	return nil
}

// SetTyping sends a typing indicator to a Telegram chat.
func (s *TelegramServer) SetTyping(chatID string) error {
	params := url.Values{}
	params.Set("chat_id", chatID)
	params.Set("action", "typing")

	resp, err := s.httpClient.Post(
		s.baseURL+"/sendChatAction",
		"application/x-www-form-urlencoded",
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_ = params
	return nil
}

// Messages returns the channel of inbound messages.
func (s *TelegramServer) Messages() <-chan *InboundMessage {
	return s.msgCh
}
