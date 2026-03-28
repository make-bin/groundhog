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

// signalCLIBase is the base URL of the signal-cli REST API sidecar.
const signalCLIBase = "http://localhost:8080"

// InboundMessage is a simplified inbound message for the gateway.
type InboundMessage struct {
	MessageID  string
	ChannelID  string
	AccountID  string
	Content    string
	ReceivedAt time.Time
}

// SignalMessage represents a message returned by the signal-cli REST API.
type SignalMessage struct {
	Envelope struct {
		Source      string `json:"source"`
		Timestamp   int64  `json:"timestamp"`
		DataMessage *struct {
			Message   string `json:"message"`
			Timestamp int64  `json:"timestamp"`
		} `json:"dataMessage,omitempty"`
	} `json:"envelope"`
}

// SignalServer implements the Signal channel plugin via signal-cli REST API.
type SignalServer struct {
	account    string
	msgCh      chan *InboundMessage
	mu         sync.Mutex
	httpClient *http.Client
}

// NewSignalServer creates a new SignalServer.
// account is the registered Signal phone number (e.g. "+15551234567").
func NewSignalServer(account string) *SignalServer {
	return &SignalServer{
		account:    account,
		msgCh:      make(chan *InboundMessage, 100),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Connect verifies that the signal-cli sidecar is reachable.
func (s *SignalServer) Connect(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		signalCLIBase+"/v1/accounts", nil)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("signal-cli sidecar unreachable: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// Start begins polling the signal-cli REST API for received messages.
func (s *SignalServer) Start(ctx context.Context) error {
	go s.poll(ctx)
	return nil
}

// poll continuously polls the signal-cli /v1/receive endpoint.
func (s *SignalServer) poll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		messages, err := s.receive(ctx)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, m := range messages {
			if m.Envelope.DataMessage == nil || m.Envelope.DataMessage.Message == "" {
				continue
			}
			msg := &InboundMessage{
				MessageID:  fmt.Sprintf("%d", m.Envelope.DataMessage.Timestamp),
				ChannelID:  "signal",
				AccountID:  m.Envelope.Source,
				Content:    m.Envelope.DataMessage.Message,
				ReceivedAt: time.Unix(m.Envelope.Timestamp/1000, 0),
			}
			select {
			case s.msgCh <- msg:
			default:
			}
		}

		time.Sleep(2 * time.Second)
	}
}

// receive calls the signal-cli REST API to fetch pending messages.
func (s *SignalServer) receive(ctx context.Context) ([]*SignalMessage, error) {
	url := fmt.Sprintf("%s/v1/receive/%s", signalCLIBase, s.account)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var messages []*SignalMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// SendMessage sends a text message via the signal-cli REST API.
func (s *SignalServer) SendMessage(recipient, text string) error {
	payload := map[string]any{
		"message":    text,
		"number":     s.account,
		"recipients": []string{recipient},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost,
		signalCLIBase+"/v2/send",
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// SetTyping sends a typing indicator to a Signal recipient.
func (s *SignalServer) SetTyping(recipient string) error {
	payload := map[string]any{
		"recipient": recipient,
		"account":   s.account,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut,
		signalCLIBase+"/v1/typing-indicator/"+s.account,
		bytes.NewReader(body))
	if err != nil {
		return err
	}
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
func (s *SignalServer) Messages() <-chan *InboundMessage {
	return s.msgCh
}

func main() {}
