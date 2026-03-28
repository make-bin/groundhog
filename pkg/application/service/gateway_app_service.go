package service

import (
	"context"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/application/eventbus"
	"github.com/make-bin/groundhog/pkg/application/hook"
	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/event"
	"github.com/make-bin/groundhog/pkg/domain/messaging/repository"
	msgservice "github.com/make-bin/groundhog/pkg/domain/messaging/service"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

// GatewayAppService orchestrates inbound message handling.
type GatewayAppService interface {
	HandleInboundMessage(ctx context.Context, req *dto.InboundMessageRequest) error
}

type gatewayAppService struct {
	MessageRepo    repository.MessageRepository `inject:""`
	RoutingService msgservice.RoutingService    `inject:""`
	EventBus       eventbus.EventBus            `inject:""`
	HookRegistry   *hook.HookRegistry           `inject:""`
}

// NewGatewayAppService creates a new GatewayAppService. Dependencies are injected via struct tags.
func NewGatewayAppService() GatewayAppService {
	return &gatewayAppService{}
}

func (s *gatewayAppService) HandleInboundMessage(ctx context.Context, req *dto.InboundMessageRequest) error {
	// Fire before-receive hook
	if s.HookRegistry != nil {
		if err := s.HookRegistry.Execute(ctx, hook.HookBeforeMessageReceive, req); err != nil {
			return fmt.Errorf("gateway_app_service: before_message_receive hook: %w", err)
		}
	}

	// Build value objects
	msgID, err := vo.NewMessageID(fmt.Sprintf("msg-%d", time.Now().UnixNano()))
	if err != nil {
		return err
	}
	channelID, err := vo.NewChannelID(req.ChannelID)
	if err != nil {
		return err
	}
	accountID, err := vo.NewAccountID(req.AccountID)
	if err != nil {
		return err
	}
	content := vo.NewMessageContent(req.Content, nil, false, nil)

	// Create aggregate
	msg := inbound_message.NewInboundMessage(msgID, channelID, accountID, content)

	// Persist
	if err := s.MessageRepo.Create(ctx, msg); err != nil {
		return fmt.Errorf("gateway_app_service: create message: %w", err)
	}

	// Resolve routing
	sessionID, err := s.RoutingService.Resolve(msg)
	if err != nil {
		return fmt.Errorf("gateway_app_service: routing: %w", err)
	}

	// Route message
	if err := msg.RouteTo(sessionID); err != nil {
		return fmt.Errorf("gateway_app_service: route: %w", err)
	}

	// Update persisted state
	if err := s.MessageRepo.Update(ctx, msg); err != nil {
		return fmt.Errorf("gateway_app_service: update message: %w", err)
	}

	// Publish domain event
	s.EventBus.Publish(event.NewMessageReceived(msgID, channelID, accountID))

	// Fire after-receive hook
	if s.HookRegistry != nil {
		if err := s.HookRegistry.Execute(ctx, hook.HookAfterMessageReceive, req); err != nil {
			return fmt.Errorf("gateway_app_service: after_message_receive hook: %w", err)
		}
	}

	return nil
}
