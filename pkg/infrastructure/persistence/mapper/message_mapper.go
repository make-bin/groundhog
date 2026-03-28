package mapper

import (
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// DomainToMessagePO converts an InboundMessage aggregate to a MessagePO.
func DomainToMessagePO(msg *inbound_message.InboundMessage) (*po.MessagePO, error) {
	return &po.MessagePO{
		MessageID:  msg.ID().Value(),
		ChannelID:  msg.ChannelID().Value(),
		AccountID:  msg.AccountID().Value(),
		Content:    msg.Content().Text(),
		Status:     int(msg.Status()),
		RoutedTo:   msg.RoutedTo(),
		ReceivedAt: msg.ReceivedAt(),
	}, nil
}

// MessagePOToDomain converts a MessagePO to an InboundMessage aggregate.
func MessagePOToDomain(p *po.MessagePO) (*inbound_message.InboundMessage, error) {
	msgID, err := vo.NewMessageID(p.MessageID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct message_id: %w", err)
	}
	channelID, err := vo.NewChannelID(p.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct channel_id: %w", err)
	}
	accountID, err := vo.NewAccountID(p.AccountID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct account_id: %w", err)
	}
	content := vo.NewMessageContent(p.Content, nil, false, nil)
	return inbound_message.ReconstructInboundMessage(
		msgID, channelID, accountID, content,
		p.ReceivedAt, p.RoutedTo, vo.MessageStatus(p.Status), nil,
	), nil
}
