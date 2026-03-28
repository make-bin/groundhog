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

const discordAPIBase = "https://discord.com/api/v10"
const discordMessageLimit = 2000

// InboundMessage is a simplified inbound message for the gateway.
type InboundMessage struct {
	MessageID  string
	ChannelID  string
	AccountID  string
	Content    string
	ReceivedAt time.Time
}

// DiscordMessage represents a Discord message object.
type DiscordMessage struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	Author    struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
	Timestamp string `json:"timestamp"`
}

// DiscordServer implements the Discord channel plugin.
type DiscordServer struct {
	token      string
	msgCh      chan *InboundMessage
	mu         sync.Mutex
	httpClient *http.Client
}

// NewDiscordServer creates a new DiscordServer.
func NewDiscordServer(token string) *DiscordServer {
	return &DiscordServer{
		token:      token,
		msgCh:      make(chan *InboundMessage, 100),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Connect establishes the Discord Gateway WebSocket connection (stub).
func (s *DiscordServer) Connect(ctx context.Context) error {
	// TODO: implement Discord Gateway WebSocket connection
	// For now, this is a stub that returns immediately
	go func() {
		<-ctx.Done()
	}()
	return nil
}

// SendMessage sends a message to a Discord channel, auto-chunking at 2000 chars.
func (s *DiscordServer) SendMessage(channelID, content string) error {
	chunks := chunkMessage(content, discordMessageLimit)
	for _, chunk := range chunks {
		if err := s.sendChunk(channelID, chunk); err != nil {
			return err
		}
	}
	return nil
}

// sendChunk sends a single message chunk to Discord.
func (s *DiscordServer) sendChunk(channelID, content string) error {
	payload := map[string]string{"content": content}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/channels/%s/messages", discordAPIBase, channelID),
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// SendEmbed sends an embed message to a Discord channel.
func (s *DiscordServer) SendEmbed(channelID, title, description string) error {
	payload := map[string]interface{}{
		"embeds": []map[string]string{
			{"title": title, "description": description},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/channels/%s/messages", discordAPIBase, channelID),
		bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bot "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

// chunkMessage splits content into chunks of at most maxLen runes.
func chunkMessage(content string, maxLen int) []string {
	runes := []rune(content)
	if len(runes) <= maxLen {
		return []string{content}
	}
	var chunks []string
	for start := 0; start < len(runes); start += maxLen {
		end := start + maxLen
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

// Messages returns the channel of inbound messages.
func (s *DiscordServer) Messages() <-chan *InboundMessage {
	return s.msgCh
}
