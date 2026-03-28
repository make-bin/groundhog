// @AI_GENERATED
package inbound_message

import (
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// InboundMessage is the aggregate root for a received message's lifecycle.
type InboundMessage struct {
	id         vo.MessageID
	channelID  vo.ChannelID
	accountID  vo.AccountID
	content    vo.MessageContent
	receivedAt time.Time
	routedTo   string
	status     vo.MessageStatus
	chunks     []vo.MessageChunk
}

// NewInboundMessage creates a new InboundMessage in Pending status.
func NewInboundMessage(id vo.MessageID, channelID vo.ChannelID, accountID vo.AccountID, content vo.MessageContent) *InboundMessage {
	return &InboundMessage{
		id:         id,
		channelID:  channelID,
		accountID:  accountID,
		content:    content,
		receivedAt: time.Now(),
		status:     vo.MessageStatusPending,
		chunks:     []vo.MessageChunk{},
	}
}

// ReconstructInboundMessage reconstructs an InboundMessage from persisted data.
// This should only be used by repository implementations.
func ReconstructInboundMessage(
	id vo.MessageID,
	channelID vo.ChannelID,
	accountID vo.AccountID,
	content vo.MessageContent,
	receivedAt time.Time,
	routedTo string,
	status vo.MessageStatus,
	chunks []vo.MessageChunk,
) *InboundMessage {
	return &InboundMessage{
		id:         id,
		channelID:  channelID,
		accountID:  accountID,
		content:    content,
		receivedAt: receivedAt,
		routedTo:   routedTo,
		status:     status,
		chunks:     chunks,
	}
}

// ID returns the message identifier.
func (m *InboundMessage) ID() vo.MessageID { return m.id }

// ChannelID returns the channel identifier.
func (m *InboundMessage) ChannelID() vo.ChannelID { return m.channelID }

// AccountID returns the account identifier.
func (m *InboundMessage) AccountID() vo.AccountID { return m.accountID }

// Content returns the message content.
func (m *InboundMessage) Content() vo.MessageContent { return m.content }

// ReceivedAt returns the time the message was received.
func (m *InboundMessage) ReceivedAt() time.Time { return m.receivedAt }

// RoutedTo returns the session ID this message was routed to.
func (m *InboundMessage) RoutedTo() string { return m.routedTo }

// Status returns the current message status.
func (m *InboundMessage) Status() vo.MessageStatus { return m.status }

// Chunks returns a copy of the message chunks.
func (m *InboundMessage) Chunks() []vo.MessageChunk {
	result := make([]vo.MessageChunk, len(m.chunks))
	copy(result, m.chunks)
	return result
}

// RouteTo routes the message to the given session. Only allowed when status is Pending.
func (m *InboundMessage) RouteTo(sessionID string) error {
	if m.status != vo.MessageStatusPending {
		return fmt.Errorf("cannot route message in status %s, must be Pending", m.status)
	}
	m.routedTo = sessionID
	m.status = vo.MessageStatusRouted
	return nil
}

// AddChunks appends chunks to the message.
func (m *InboundMessage) AddChunks(chunks []vo.MessageChunk) {
	m.chunks = append(m.chunks, chunks...)
}

// MarkDelivered sets the message status to Delivered.
func (m *InboundMessage) MarkDelivered() {
	m.status = vo.MessageStatusDelivered
}

// MarkFailed sets the message status to Failed.
func (m *InboundMessage) MarkFailed() {
	m.status = vo.MessageStatusFailed
}

// @AI_GENERATED: end
