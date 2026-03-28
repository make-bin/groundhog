package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const slackAPIBase = "https://slack.com/api"

// InboundMessage is a simplified inbound message for the gateway.
type InboundMessage struct {
	MessageID  string
	ChannelID  string
	AccountID  string
	Content    string
	ReceivedAt time.Time
}

// SlackEvent represents a Slack Events API event envelope.
type SlackEvent struct {
	Type    string      `json:"type"`
	EventID string      `json:"event_id"`
	Event   *SlackInner `json:"event,omitempty"`
}

// SlackInner represents the inner event payload.
type SlackInner struct {
	Type        string `json:"type"`
	ClientMsgID string `json:"client_msg_id,omitempty"`
	Text        string `json:"text"`
	User        string `json:"user"`
	Channel     string `json:"channel"`
	Timestamp   string `json:"ts"`
}

// SlackServer implements the Slack channel plugin.
type SlackServer struct {
	token      string
	msgCh      chan *InboundMessage
	mu         sync.Mutex
	httpClient *http.Client
	cursor     string
}

// NewSlackServer creates a new SlackServer with the given bot token.
func NewSlackServer(token string) *SlackServer {
	return &SlackServer{
		token:      token,
		msgCh:      make(chan *InboundMessage, 100),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Connect initialises the Slack connection (stub — real impl would open a
// Socket Mode WebSocket or register an Events API endpoint).
func (s *SlackServer) Connect(ctx context.Context) error {
	go func() {
		<-ctx.Done()
	}()
	return nil
}

// Start begins polling for new messages via conversations.history.
func (s *SlackServer) Start(ctx context.Context, channelID string) error {
	go s.poll(ctx, channelID)
	return nil
}

// poll continuously polls the Slack conversations.history endpoint.
func (s *SlackServer) poll(ctx context.Context, channelID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		messages, nextCursor, err := s.fetchMessages(ctx, channelID)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		for _, m := range messages {
			msg := &InboundMessage{
				MessageID:  m.Timestamp,
				ChannelID:  channelID,
				AccountID:  m.User,
				Content:    m.Text,
				ReceivedAt: time.Now(),
			}
			select {
			case s.msgCh <- msg:
			default:
			}
		}

		s.mu.Lock()
		s.cursor = nextCursor
		s.mu.Unlock()

		time.Sleep(2 * time.Second)
	}
}

// fetchMessages calls conversations.history and returns new messages.
func (s *SlackServer) fetchMessages(ctx context.Context, channelID string) ([]*SlackInner, string, error) {
	s.mu.Lock()
	cursor := s.cursor
	s.mu.Unlock()

	params := url.Values{}
	params.Set("channel", channelID)
	params.Set("limit", "20")
	if cursor != "" {
		params.Set("cursor", cursor)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		slackAPIBase+"/conversations.history?"+params.Encode(), nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		OK               bool          `json:"ok"`
		Messages         []*SlackInner `json:"messages"`
		ResponseMetadata struct {
			NextCursor string `json:"next_cursor"`
		} `json:"response_metadata"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}
	if !result.OK {
		return nil, "", fmt.Errorf("slack API error")
	}
	return result.Messages, result.ResponseMetadata.NextCursor, nil
}

// SendMessage sends a text message to a Slack channel via chat.postMessage.
func (s *SlackServer) SendMessage(channelID, text string) error {
	payload := map[string]string{
		"channel": channelID,
		"text":    text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost,
		slackAPIBase+"/chat.postMessage",
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// SetTyping sends a typing indicator to a Slack channel.
// Slack does not expose a public typing API for bots; we send an ephemeral
// "typing…" message as a best-effort indicator instead.
func (s *SlackServer) SetTyping(channelID string) error {
	payload := map[string]any{
		"channel": channelID,
		"text":    "_typing…_",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost,
		slackAPIBase+"/chat.postMessage",
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
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
func (s *SlackServer) Messages() <-chan *InboundMessage {
	return s.msgCh
}

func main() {}
