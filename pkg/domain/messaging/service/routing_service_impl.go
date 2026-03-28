package service

import (
	"fmt"
	"sync"

	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/vo"
)

type routingServiceImpl struct {
	mu       sync.RWMutex
	bindings map[string]string // key: "channelID:accountID", value: sessionID
}

// NewRoutingService creates a new RoutingService implementation.
func NewRoutingService() RoutingService {
	return &routingServiceImpl{
		bindings: make(map[string]string),
	}
}

func (s *routingServiceImpl) Resolve(msg *inbound_message.InboundMessage) (string, error) {
	key := fmt.Sprintf("%s:%s", msg.ChannelID().Value(), msg.AccountID().Value())
	s.mu.RLock()
	sessionID, ok := s.bindings[key]
	s.mu.RUnlock()
	if ok {
		return sessionID, nil
	}
	// No existing binding — create a new session ID
	newSessionID := fmt.Sprintf("sess-%s-%s", msg.ChannelID().Value(), msg.AccountID().Value())
	s.mu.Lock()
	s.bindings[key] = newSessionID
	s.mu.Unlock()
	return newSessionID, nil
}

// BindAccount creates or updates an account binding.
func (s *routingServiceImpl) BindAccount(binding vo.AccountBinding) {
	key := fmt.Sprintf("%s:%s", binding.ChannelID().Value(), binding.AccountID().Value())
	s.mu.Lock()
	s.bindings[key] = binding.SessionID()
	s.mu.Unlock()
}
