package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/make-bin/groundhog/pkg/domain/messaging/aggregate/inbound_message"
	"github.com/make-bin/groundhog/pkg/domain/messaging/service"
	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// SessionCreator is the minimal interface needed to create a session for a given agent.
// Implemented by AgentAppService; defined here to avoid import cycles.
type SessionCreator interface {
	GetOrCreateSession(ctx context.Context, agentID, userID, channelType string) (string, error)
}

// AgentRoutingService implements RoutingService with openclaw-style binding rules.
// It resolves inbound messages to sessions by:
//  1. Matching bindings (channel_id / channel type / account_id) → agent_id
//  2. Falling back to the default agent
//  3. Calling SessionCreator.GetOrCreateSession to get/create the target session
type AgentRoutingService struct {
	registry       *Registry
	bindings       []config.BindingConfig
	sessionCreator SessionCreator
	log            logger.Logger

	// in-memory cache: "agentID:userID" → sessionID (avoids repeated DB lookups within a process)
	mu    sync.RWMutex
	cache map[string]string
}

// NewAgentRoutingService creates an AgentRoutingService.
// sessionCreator is wired after DI population to avoid circular dependencies.
func NewAgentRoutingService(registry *Registry, bindings []config.BindingConfig, log logger.Logger) *AgentRoutingService {
	return &AgentRoutingService{
		registry: registry,
		bindings: bindings,
		log:      log,
		cache:    make(map[string]string),
	}
}

// SetSessionCreator wires the session creator after DI population.
func (s *AgentRoutingService) SetSessionCreator(sc SessionCreator) {
	s.sessionCreator = sc
}

// Resolve implements RoutingService.
// It matches the message against bindings to find the target agent,
// then gets or creates a session for that agent+user combination.
func (s *AgentRoutingService) Resolve(msg *inbound_message.InboundMessage) (string, error) {
	agentID := s.resolveAgentID(msg)
	userID := msg.AccountID().Value()

	cacheKey := agentID + ":" + userID
	s.mu.RLock()
	if sid, ok := s.cache[cacheKey]; ok {
		s.mu.RUnlock()
		return sid, nil
	}
	s.mu.RUnlock()

	if s.sessionCreator == nil {
		// Fallback: synthesize a deterministic session key (no DB)
		return fmt.Sprintf("agent:%s:%s", agentID, userID), nil
	}

	channelType := s.channelTypeFromID(msg.ChannelID().Value())
	sessionID, err := s.sessionCreator.GetOrCreateSession(context.Background(), agentID, userID, channelType)
	if err != nil {
		return "", fmt.Errorf("agent_routing: get_or_create_session: %w", err)
	}

	s.mu.Lock()
	s.cache[cacheKey] = sessionID
	s.mu.Unlock()

	s.log.Info("routed message to session",
		"agent_id", agentID,
		"user_id", userID,
		"session_id", sessionID,
		"channel_id", msg.ChannelID().Value(),
	)
	return sessionID, nil
}

// resolveAgentID finds the best matching agent for the message.
// Binding priority (most specific first):
//  1. channel_id + account_id match
//  2. channel_id match only
//  3. channel type + account_id match
//  4. channel type match only
//  5. default agent
func (s *AgentRoutingService) resolveAgentID(msg *inbound_message.InboundMessage) string {
	channelID := msg.ChannelID().Value()
	accountID := msg.AccountID().Value()
	channelType := s.channelTypeFromID(channelID)

	type candidate struct {
		agentID    string
		specificity int // higher = more specific
	}
	var best candidate

	for _, b := range s.bindings {
		agentID := NormalizeAgentID(b.AgentID)
		if agentID == "" {
			continue
		}
		spec := s.matchSpecificity(b.Match, channelID, channelType, accountID)
		if spec < 0 {
			continue // no match
		}
		if spec > best.specificity || best.agentID == "" {
			best = candidate{agentID: agentID, specificity: spec}
		}
	}

	if best.agentID != "" {
		return best.agentID
	}
	return s.registry.DefaultID()
}

// matchSpecificity returns how specifically a BindingMatchConfig matches the message.
// Returns -1 if no match. Higher values = more specific.
func (s *AgentRoutingService) matchSpecificity(
	match config.BindingMatchConfig,
	channelID, channelType, accountID string,
) int {
	// channel_id match (most specific)
	if match.ChannelID != "" {
		if match.ChannelID != channelID {
			return -1
		}
		if match.AccountID != "" {
			if match.AccountID != accountID {
				return -1
			}
			return 4 // channel_id + account_id
		}
		return 3 // channel_id only
	}

	// channel type match
	if match.Channel != "" {
		if !strings.EqualFold(match.Channel, channelType) {
			return -1
		}
		if match.AccountID != "" {
			if match.AccountID != accountID {
				return -1
			}
			return 2 // channel type + account_id
		}
		return 1 // channel type only
	}

	// account_id only (no channel constraint)
	if match.AccountID != "" {
		if match.AccountID != accountID {
			return -1
		}
		return 1
	}

	// empty match = catch-all (lowest specificity)
	return 0
}

// channelTypeFromID extracts the channel type from a channel ID.
// Channel IDs are typically prefixed: "discord-guild-123" → "discord".
func (s *AgentRoutingService) channelTypeFromID(channelID string) string {
	parts := strings.SplitN(channelID, "-", 2)
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}
	return channelID
}

// Ensure AgentRoutingService implements RoutingService.
var _ service.RoutingService = (*AgentRoutingService)(nil)
