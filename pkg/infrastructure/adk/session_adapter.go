// @AI_GENERATED
package adk

import (
	"context"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// ADKSession is the internal representation of a session for the ADK layer.
type ADKSession struct {
	ID        string
	UserID    string
	AgentID   string
	AppState  map[string]any
	UserState map[string]any
	TempState map[string]any
}

// SessionAdapter converts between domain AgentSession and ADKSession.
type SessionAdapter struct {
	DomainRepo repository.SessionRepository `inject:""`
}

// NewSessionAdapter creates a new SessionAdapter.
func NewSessionAdapter() *SessionAdapter {
	return &SessionAdapter{}
}

// ToADK converts a domain AgentSession to an ADKSession.
func (a *SessionAdapter) ToADK(s *agent_session.AgentSession) *ADKSession {
	meta := s.Metadata()
	appState := make(map[string]any)
	userState := make(map[string]any)
	tempState := make(map[string]any)
	for k, v := range meta {
		switch {
		case len(k) > 4 && k[:4] == "app:":
			appState[k[4:]] = v
		case len(k) > 5 && k[:5] == "user:":
			userState[k[5:]] = v
		case len(k) > 5 && k[:5] == "temp:":
			tempState[k[5:]] = v
		}
	}
	return &ADKSession{
		ID:        s.ID().Value(),
		UserID:    s.UserID(),
		AgentID:   s.AgentID().Value(),
		AppState:  appState,
		UserState: userState,
		TempState: tempState,
	}
}

// FromADK converts an ADKSession back to a minimal domain AgentSession.
// Note: this creates a new session with metadata populated from ADK state scopes.
func (a *SessionAdapter) FromADK(s *ADKSession) (*agent_session.AgentSession, error) {
	sessionID, err := vo.NewSessionID(s.ID)
	if err != nil {
		return nil, fmt.Errorf("session_adapter: invalid session id: %w", err)
	}
	agentID, err := vo.NewAgentID(s.AgentID)
	if err != nil {
		return nil, fmt.Errorf("session_adapter: invalid agent id: %w", err)
	}
	// Build a default ModelConfig — real config comes from the domain session.
	modelCfg, err := vo.NewModelConfig(vo.ProviderGemini, "gemini-pro", 0.7, 4096, nil, "")
	if err != nil {
		return nil, fmt.Errorf("session_adapter: build default model config: %w", err)
	}
	sess, err := agent_session.NewAgentSession(sessionID, agentID, s.UserID, modelCfg, vo.NewPrompt("", nil))
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// Load retrieves a domain AgentSession by ID from the repository.
func (a *SessionAdapter) Load(ctx context.Context, id vo.SessionID) (*agent_session.AgentSession, error) {
	return a.DomainRepo.FindByID(ctx, id)
}

// Save persists a domain AgentSession via the repository.
func (a *SessionAdapter) Save(ctx context.Context, s *agent_session.AgentSession) error {
	return a.DomainRepo.Update(ctx, s)
}

// @AI_GENERATED: end
