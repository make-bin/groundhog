package event

import (
	"time"

	"github.com/make-bin/groundhog/pkg/domain/messaging/entity"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// DomainEvent is the interface all messaging domain events must implement.
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// MessageReceived is raised when a new inbound message is received from a channel.
type MessageReceived struct {
	MessageID  vo.MessageID
	ChannelID  vo.ChannelID
	AccountID  vo.AccountID
	occurredAt time.Time
}

// NewMessageReceived creates a new MessageReceived event.
func NewMessageReceived(messageID vo.MessageID, channelID vo.ChannelID, accountID vo.AccountID) MessageReceived {
	return MessageReceived{
		MessageID:  messageID,
		ChannelID:  channelID,
		AccountID:  accountID,
		occurredAt: time.Now(),
	}
}

func (e MessageReceived) OccurredAt() time.Time { return e.occurredAt }
func (e MessageReceived) EventType() string     { return "messaging.message.received" }

// MessageRouted is raised when a message has been successfully routed to a session.
type MessageRouted struct {
	MessageID  vo.MessageID
	SessionID  string
	occurredAt time.Time
}

// NewMessageRouted creates a new MessageRouted event.
func NewMessageRouted(messageID vo.MessageID, sessionID string) MessageRouted {
	return MessageRouted{
		MessageID:  messageID,
		SessionID:  sessionID,
		occurredAt: time.Now(),
	}
}

func (e MessageRouted) OccurredAt() time.Time { return e.occurredAt }
func (e MessageRouted) EventType() string     { return "messaging.message.routed" }

// MessageDeliveryFailed is raised when message delivery to a channel fails.
type MessageDeliveryFailed struct {
	MessageID  vo.MessageID
	Reason     string
	occurredAt time.Time
}

// NewMessageDeliveryFailed creates a new MessageDeliveryFailed event.
func NewMessageDeliveryFailed(messageID vo.MessageID, reason string) MessageDeliveryFailed {
	return MessageDeliveryFailed{
		MessageID:  messageID,
		Reason:     reason,
		occurredAt: time.Now(),
	}
}

func (e MessageDeliveryFailed) OccurredAt() time.Time { return e.occurredAt }
func (e MessageDeliveryFailed) EventType() string     { return "messaging.message.delivery_failed" }

// ChannelStatusChanged is raised when a channel transitions between operational states.
type ChannelStatusChanged struct {
	ChannelID  vo.ChannelID
	OldStatus  entity.ChannelStatus
	NewStatus  entity.ChannelStatus
	occurredAt time.Time
}

// NewChannelStatusChanged creates a new ChannelStatusChanged event.
func NewChannelStatusChanged(channelID vo.ChannelID, oldStatus, newStatus entity.ChannelStatus) ChannelStatusChanged {
	return ChannelStatusChanged{
		ChannelID:  channelID,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		occurredAt: time.Now(),
	}
}

func (e ChannelStatusChanged) OccurredAt() time.Time { return e.occurredAt }
func (e ChannelStatusChanged) EventType() string     { return "messaging.channel.status_changed" }
