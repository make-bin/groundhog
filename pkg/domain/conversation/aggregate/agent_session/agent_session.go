package agent_session

import (
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// AgentSession is the aggregate root for a conversation session.
type AgentSession struct {
	id              vo.SessionID
	agentID         vo.AgentID
	userID          string
	turns           []entity.Turn
	activeModel     vo.ModelConfig
	tools           []entity.ToolDefinition
	systemPrompt    vo.Prompt
	skills          []string
	state           vo.SessionState
	createdAt       time.Time
	lastActiveAt    time.Time
	metadata        map[string]any
	compactionCount int
}

// NewAgentSession creates a new AgentSession. Returns an error if userID is empty.
func NewAgentSession(id vo.SessionID, agentID vo.AgentID, userID string, activeModel vo.ModelConfig, systemPrompt vo.Prompt) (*AgentSession, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID must not be empty")
	}
	now := time.Now()
	return &AgentSession{
		id:           id,
		agentID:      agentID,
		userID:       userID,
		turns:        []entity.Turn{},
		activeModel:  activeModel,
		tools:        []entity.ToolDefinition{},
		systemPrompt: systemPrompt,
		skills:       []string{},
		state:        vo.SessionStateActive,
		createdAt:    now,
		lastActiveAt: now,
		metadata:     map[string]any{},
	}, nil
}

// ID returns the session identifier.
func (s *AgentSession) ID() vo.SessionID { return s.id }

// AgentID returns the agent identifier.
func (s *AgentSession) AgentID() vo.AgentID { return s.agentID }

// UserID returns the user identifier.
func (s *AgentSession) UserID() string { return s.userID }

// Turns returns a copy of the conversation turns.
func (s *AgentSession) Turns() []entity.Turn {
	result := make([]entity.Turn, len(s.turns))
	copy(result, s.turns)
	return result
}

// ActiveModel returns the active model configuration.
func (s *AgentSession) ActiveModel() vo.ModelConfig { return s.activeModel }

// Tools returns a copy of the tool definitions.
func (s *AgentSession) Tools() []entity.ToolDefinition {
	result := make([]entity.ToolDefinition, len(s.tools))
	copy(result, s.tools)
	return result
}

// SystemPrompt returns the system prompt.
func (s *AgentSession) SystemPrompt() vo.Prompt { return s.systemPrompt }

// Skills returns a copy of the skill identifiers.
func (s *AgentSession) Skills() []string {
	result := make([]string, len(s.skills))
	copy(result, s.skills)
	return result
}

// State returns the current session state.
func (s *AgentSession) State() vo.SessionState { return s.state }

// CreatedAt returns the time the session was created.
func (s *AgentSession) CreatedAt() time.Time { return s.createdAt }

// LastActiveAt returns the time of the last activity.
func (s *AgentSession) LastActiveAt() time.Time { return s.lastActiveAt }

// Metadata returns the session metadata map.
func (s *AgentSession) Metadata() map[string]any { return s.metadata }

// AddTurn appends a turn and updates lastActiveAt.
func (s *AgentSession) AddTurn(turn entity.Turn) {
	s.turns = append(s.turns, turn)
	s.lastActiveAt = time.Now()
}

// NeedsCompaction returns true when len(turns) > threshold.
func (s *AgentSession) NeedsCompaction(threshold int) bool {
	return len(s.turns) > threshold
}

// CompactTurns replaces all but the most recent keepRecent turns with a single summary Turn.
// If len(turns) <= keepRecent, this is a no-op.
func (s *AgentSession) CompactTurns(summary string, keepRecent int) {
	if len(s.turns) <= keepRecent {
		return
	}
	recent := s.turns[len(s.turns)-keepRecent:]
	summaryTurn := entity.NewSummaryTurn(
		fmt.Sprintf("summary-%d", time.Now().UnixNano()),
		summary,
	)
	newTurns := make([]entity.Turn, 0, 1+keepRecent)
	newTurns = append(newTurns, *summaryTurn)
	newTurns = append(newTurns, recent...)
	s.turns = newTurns
	s.compactionCount++
}

// CompactionCount returns the number of times compaction has been performed on this session.
func (s *AgentSession) CompactionCount() int { return s.compactionCount }

// Archive transitions the session to Archived. Only Active sessions can be archived.
func (s *AgentSession) Archive() error {
	if s.state != vo.SessionStateActive {
		return fmt.Errorf("only active sessions can be archived, current state: %s", s.state)
	}
	s.state = vo.SessionStateArchived
	return nil
}

// UpdateModel replaces the active model configuration.
func (s *AgentSession) UpdateModel(cfg vo.ModelConfig) {
	s.activeModel = cfg
}

// SetTools replaces the session tools with a copy of the provided slice.
func (s *AgentSession) SetTools(tools []entity.ToolDefinition) {
	result := make([]entity.ToolDefinition, len(tools))
	copy(result, tools)
	s.tools = result
}

// SetSkills replaces the session skill identifiers.
func (s *AgentSession) SetSkills(skills []string) {
	result := make([]string, len(skills))
	copy(result, skills)
	s.skills = result
}

// ReconstructAgentSession reconstructs an AgentSession from persisted data.
// This should only be used by repository implementations.
func ReconstructAgentSession(
	id vo.SessionID,
	agentID vo.AgentID,
	userID string,
	turns []entity.Turn,
	activeModel vo.ModelConfig,
	tools []entity.ToolDefinition,
	systemPrompt vo.Prompt,
	skills []string,
	state vo.SessionState,
	createdAt time.Time,
	lastActiveAt time.Time,
	metadata map[string]any,
	compactionCount int,
) *AgentSession {
	return &AgentSession{
		id:              id,
		agentID:         agentID,
		userID:          userID,
		turns:           turns,
		activeModel:     activeModel,
		tools:           tools,
		systemPrompt:    systemPrompt,
		skills:          skills,
		state:           state,
		createdAt:       createdAt,
		lastActiveAt:    lastActiveAt,
		metadata:        metadata,
		compactionCount: compactionCount,
	}
}
