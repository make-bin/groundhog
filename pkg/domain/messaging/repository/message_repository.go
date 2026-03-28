package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// MessageFilter is a type-safe filter for querying messages.
type MessageFilter struct {
	ChannelID *vo.ChannelID
	AccountID *vo.AccountID
	Status    *vo.MessageStatus
}

// MessageRepository defines the data access contract for the InboundMessage aggregate.
type MessageRepository interface {
	Create(ctx context.Context, msg *inbound_message.InboundMessage) error
	FindByID(ctx context.Context, id vo.MessageID) (*inbound_message.InboundMessage, error)
	Update(ctx context.Context, msg *inbound_message.InboundMessage) error
	List(ctx context.Context, filter MessageFilter, offset, limit int) ([]*inbound_message.InboundMessage, int, error)
}
