package repository

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// SessionFilter is a type-safe filter struct for querying sessions.
type SessionFilter struct {
	UserID  *string
	AgentID *vo.AgentID
	State   *vo.SessionState
}

// SessionRepository defines the data access contract for the AgentSession aggregate.
type SessionRepository interface {
	// Create persists a new AgentSession aggregate.
	Create(ctx context.Context, session *agent_session.AgentSession) error

	// FindByID retrieves a complete AgentSession aggregate by its ID.
	FindByID(ctx context.Context, id vo.SessionID) (*agent_session.AgentSession, error)

	// Update persists changes to an existing AgentSession aggregate.
	Update(ctx context.Context, session *agent_session.AgentSession) error

	// Delete removes an AgentSession aggregate by its ID.
	Delete(ctx context.Context, id vo.SessionID) error

	// List retrieves AgentSession aggregates matching the filter with pagination.
	List(ctx context.Context, filter SessionFilter, offset, limit int) ([]*agent_session.AgentSession, int, error)

	// Archive transitions an AgentSession to the archived state.
	Archive(ctx context.Context, id vo.SessionID) error
}
