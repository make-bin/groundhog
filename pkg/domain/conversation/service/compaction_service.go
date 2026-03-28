package service

import (
	"context"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
)

// CompactionService encapsulates context compaction domain logic.
type CompactionService interface {
	// Compact performs compaction on the session, keeping the most recent keepRecent turns.
	// If LLM call fails, logs the error and returns nil (does not interrupt the caller).
	Compact(ctx context.Context, sess *agent_session.AgentSession, keepRecent int) error
}
