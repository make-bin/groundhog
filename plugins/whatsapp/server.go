package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const whatsappAPIBase = "https://graph.facebook.com/v18.0"

// InboundMessage is a simplified inbound message for the gateway.
type InboundMessage struct {
	MessageID  string
	ChannelID  string
	AccountID  string
	Content    string
	ReceivedAt time.Time
}

// WAMessage represents a WhatsApp Business Cloud API inbound message.
type WAMessage struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      struct {
		Body string `json:"body"`
	} `json:"text"`
}

// WAWebhookPayload is the top-level webhook notification envelope.
type WAWebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		Changes []struct {
			Value struct {
				Messages []*WAMessage `json:"messages"`
				Metadata struct {
					PhoneNumberID string `json:"phone_number_id"`
				} `json:"metadata"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// WhatsAppServer implements the WhatsApp channel plugin.
type WhatsAppServer struct {
	accessToken   string
	phoneNumberID string
	msgCh         chan *InboundMessage
	mu            sync.Mutex
	httpClient    *http.Client
}

// NewWhatsAppServer creates a new WhatsAppServer.
func NewWhatsAppServer(accessToken, phoneNumberID string) *WhatsAppServer {
	return &WhatsAppServer{
		accessToken:   accessToken,
		phoneNumberID: phoneNumberID,
		msgCh:         make(chan *InboundMessage, 100),
		httpClient:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Connect initialises the WhatsApp connection.
// In production this would register the webhook URL with the Graph API.
func (s *WhatsAppServer) Connect(ctx context.Context) error {
	go func() {
		<-ctx.Done()
	}()
	return nil
}

// Start begins listening for inbound webhook events.
// In production a real HTTP server would receive POST callbacks from Meta;
// this stub starts a goroutine that waits for context cancellation.
func (s *WhatsAppServer) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
	}()
	return nil
}

// HandleWebhook processes a raw webhook payload from Meta and enqueues
// any inbound text messages.
func (s *WhatsAppServer) HandleWebhook(data []byte) error {
	var payload WAWebhookPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			for _, m := range change.Value.Messages {
				if m.Type != "text" {
					continue
				}
				msg := &InboundMessage{
					MessageID:  m.ID,
					ChannelID:  change.Value.Metadata.PhoneNumberID,
					AccountID:  m.From,
					Content:    m.Text.Body,
					ReceivedAt: time.Now(),
				}
				select {
				case s.msgCh <- msg:
				default:
				}
			}
		}
	}
	return nil
}

// SendMessage sends a text message via the WhatsApp Business Cloud API.
func (s *WhatsAppServer) SendMessage(to, text string) error {
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text":              map[string]string{"body": text},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/messages", whatsappAPIBase, s.phoneNumberID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// SetTyping sends a "read" receipt which triggers the typing indicator on the
// recipient's device (WhatsApp Business API does not expose a standalone
// typing action endpoint).
func (s *WhatsAppServer) SetTyping(messageID string) error {
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/%s/messages", whatsappAPIBase, s.phoneNumberID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// Messages returns the channel of inbound messages.
func (s *WhatsAppServer) Messages() <-chan *InboundMessage {
	return s.msgCh
}

func main() {}
