package adk

import (
	"context"
	"fmt"
	"strings"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	conversation_service "github.com/make-bin/groundhog/pkg/domain/conversation/service"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// compactionServiceImpl implements conversation_service.CompactionService using the ADK layer.
type compactionServiceImpl struct {
	ModelAdapter *ModelAdapter `inject:""`
	Logger       logger.Logger `inject:"logger"`
}

// NewCompactionService creates a new CompactionService.
func NewCompactionService() conversation_service.CompactionService {
	return &compactionServiceImpl{}
}

// Compact performs context compaction on the session by summarizing turns to be compressed.
// If the LLM call fails or returns an empty summary, it logs a warning and returns nil.
func (s *compactionServiceImpl) Compact(ctx context.Context, sess *agent_session.AgentSession, keepRecent int) error {
	turns := sess.Turns()
	if len(turns) <= keepRecent {
		return nil
	}

	// Build compaction prompt from turns that will be compressed
	turnsToCompress := turns[:len(turns)-keepRecent]
	prompt := buildCompactionPrompt(turnsToCompress)

	// Get LLM for the session's active model
	llm, err := s.ModelAdapter.ToADKModel(ctx, sess.ActiveModel())
	if err != nil {
		if s.Logger != nil {
			s.Logger.Warn("compaction: failed to get LLM model", "error", err)
		}
		return nil
	}

	// Generate summary
	summary, err := llm.GenerateContent(ctx, prompt)
	if err != nil {
		if s.Logger != nil {
			s.Logger.Warn("compaction: LLM call failed, skipping compaction", "error", err)
		}
		return nil
	}

	// Skip if summary is empty
	if strings.TrimSpace(summary) == "" {
		if s.Logger != nil {
			s.Logger.Warn("compaction: LLM returned empty summary, skipping compaction")
		}
		return nil
	}

	sess.CompactTurns(summary, keepRecent)
	return nil
}

// buildCompactionPrompt constructs the prompt for the LLM to summarize the given turns.
func buildCompactionPrompt(turns []entity.Turn) string {
	var sb strings.Builder
	sb.WriteString("Please summarize the following conversation history into a concise summary that preserves the key information, decisions, and context:\n\n")
	for _, t := range turns {
		if t.IsSummary() {
			sb.WriteString(fmt.Sprintf("[Previous Summary]: %s\n\n", t.Response()))
		} else {
			if t.UserInput() != "" {
				sb.WriteString(fmt.Sprintf("User: %s\n", t.UserInput()))
			}
			if t.Response() != "" {
				sb.WriteString(fmt.Sprintf("Assistant: %s\n\n", t.Response()))
			}
		}
	}
	return sb.String()
}
