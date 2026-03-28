// @AI_GENERATED
package entity

import (
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// Turn is an entity representing a single conversation turn.
type Turn struct {
	id          string
	userInput   string
	response    string
	toolCalls   []ToolCall
	modelUsed   string
	tokenUsage  vo.TokenUsage
	startedAt   time.Time
	completedAt time.Time
	err         error
	isSummary   bool
}

// NewTurn creates a new Turn with the given id and user input.
// The startedAt field is set to the current time.
func NewTurn(id string, userInput string) *Turn {
	return &Turn{
		id:        id,
		userInput: userInput,
		startedAt: time.Now(),
	}
}

// ID returns the turn identifier.
func (t *Turn) ID() string { return t.id }

// UserInput returns the user input text.
func (t *Turn) UserInput() string { return t.userInput }

// Response returns the agent response text.
func (t *Turn) Response() string { return t.response }

// ToolCalls returns a copy of the tool calls made during this turn.
func (t *Turn) ToolCalls() []ToolCall {
	calls := make([]ToolCall, len(t.toolCalls))
	copy(calls, t.toolCalls)
	return calls
}

// ModelUsed returns the model identifier used for this turn.
func (t *Turn) ModelUsed() string { return t.modelUsed }

// TokenUsage returns the token usage for this turn.
func (t *Turn) TokenUsage() vo.TokenUsage { return t.tokenUsage }

// StartedAt returns the time the turn started.
func (t *Turn) StartedAt() time.Time { return t.startedAt }

// CompletedAt returns the time the turn completed.
func (t *Turn) CompletedAt() time.Time { return t.completedAt }

// Err returns the error that occurred during this turn, if any.
func (t *Turn) Err() error { return t.err }

// IsSummary returns true if this turn is a compaction summary rather than a real conversation turn.
func (t *Turn) IsSummary() bool { return t.isSummary }

// SetResponse sets the agent response text.
func (t *Turn) SetResponse(response string) {
	t.response = response
}

// AddToolCall appends a tool call to this turn.
func (t *Turn) AddToolCall(tc ToolCall) {
	t.toolCalls = append(t.toolCalls, tc)
}

// SetModelUsed sets the model identifier used for this turn.
func (t *Turn) SetModelUsed(model string) {
	t.modelUsed = model
}

// Complete marks the turn as completed with the given token usage.
// The completedAt field is set to the current time.
func (t *Turn) Complete(usage vo.TokenUsage) {
	t.tokenUsage = usage
	t.completedAt = time.Now()
}

// @AI_GENERATED: end

// NewSummaryTurn creates a synthetic Turn representing a compaction summary.
// userInput is empty, isSummary is true, and response holds the summary text.
func NewSummaryTurn(id string, summary string) *Turn {
	return &Turn{
		id:        id,
		userInput: "",
		response:  summary,
		isSummary: true,
		startedAt: time.Now(),
	}
}

// ReconstructTurn reconstructs a Turn from persisted data.
// This should only be used by repository implementations.
func ReconstructTurn(
	id string,
	userInput string,
	response string,
	modelUsed string,
	toolCalls []ToolCall,
	tokenUsage vo.TokenUsage,
	startedAt time.Time,
	completedAt time.Time,
	isSummary bool,
) *Turn {
	return &Turn{
		id:          id,
		userInput:   userInput,
		response:    response,
		modelUsed:   modelUsed,
		toolCalls:   toolCalls,
		tokenUsage:  tokenUsage,
		startedAt:   startedAt,
		completedAt: completedAt,
		isSummary:   isSummary,
	}
}
