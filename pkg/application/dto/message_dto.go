package dto

import "time"

// InboundMessageRequest is the request DTO for handling an inbound message from a channel.
type InboundMessageRequest struct {
	MessageID string `json:"message_id" binding:"required"`
	ChannelID string `json:"channel_id" binding:"required"`
	AccountID string `json:"account_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// ChannelDTO is the response DTO for a channel.
type ChannelDTO struct {
	ID           string   `json:"id"`
	ChannelType  string   `json:"channel_type"`
	PluginID     string   `json:"plugin_id"`
	Status       string   `json:"status"`
	Capabilities []string `json:"capabilities"`
}

// CreateChannelRequest is the request DTO for creating a channel.
type CreateChannelRequest struct {
	ChannelID   string `json:"channel_id" binding:"required"`
	ChannelType string `json:"channel_type" binding:"required"`
	PluginID    string `json:"plugin_id"`
}

// MessageDTO is the response DTO for a message.
type MessageDTO struct {
	ID         string    `json:"id"`
	ChannelID  string    `json:"channel_id"`
	AccountID  string    `json:"account_id"`
	Content    string    `json:"content"`
	Status     string    `json:"status"`
	RoutedTo   string    `json:"routed_to,omitempty"`
	ReceivedAt time.Time `json:"received_at"`
}
