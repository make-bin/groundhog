package event

import "time"

// DomainEvent is the interface all conversation domain events must implement.
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
}

// AgentTurnStarted is raised when an agent turn begins.
type AgentTurnStarted struct {
	SessionID  string
	TurnID     string
	UserInput  string
	occurredAt time.Time
}

func NewAgentTurnStarted(sessionID, turnID, userInput string) AgentTurnStarted {
	return AgentTurnStarted{SessionID: sessionID, TurnID: turnID, UserInput: userInput, occurredAt: time.Now()}
}
func (e AgentTurnStarted) OccurredAt() time.Time { return e.occurredAt }
func (e AgentTurnStarted) EventType() string     { return "agent.turn.started" }

// AgentTurnCompleted is raised when an agent turn finishes successfully.
type AgentTurnCompleted struct {
	SessionID        string
	TurnID           string
	Response         string
	PromptTokens     int
	CompletionTokens int
	occurredAt       time.Time
}

func NewAgentTurnCompleted(sessionID, turnID, response string, promptTokens, completionTokens int) AgentTurnCompleted {
	return AgentTurnCompleted{
		SessionID:        sessionID,
		TurnID:           turnID,
		Response:         response,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		occurredAt:       time.Now(),
	}
}
func (e AgentTurnCompleted) OccurredAt() time.Time { return e.occurredAt }
func (e AgentTurnCompleted) EventType() string     { return "agent.turn.completed" }

// ToolExecutionRequested is raised when a tool call is about to be executed.
type ToolExecutionRequested struct {
	SessionID  string
	TurnID     string
	ToolName   string
	Args       map[string]any
	occurredAt time.Time
}

func NewToolExecutionRequested(sessionID, turnID, toolName string, args map[string]any) ToolExecutionRequested {
	return ToolExecutionRequested{SessionID: sessionID, TurnID: turnID, ToolName: toolName, Args: args, occurredAt: time.Now()}
}
func (e ToolExecutionRequested) OccurredAt() time.Time { return e.occurredAt }
func (e ToolExecutionRequested) EventType() string     { return "agent.tool.execution_requested" }

// ToolExecutionCompleted is raised when a tool call finishes.
type ToolExecutionCompleted struct {
	SessionID  string
	TurnID     string
	ToolName   string
	Output     string
	IsError    bool
	occurredAt time.Time
}

func NewToolExecutionCompleted(sessionID, turnID, toolName, output string, isError bool) ToolExecutionCompleted {
	return ToolExecutionCompleted{SessionID: sessionID, TurnID: turnID, ToolName: toolName, Output: output, IsError: isError, occurredAt: time.Now()}
}
func (e ToolExecutionCompleted) OccurredAt() time.Time { return e.occurredAt }
func (e ToolExecutionCompleted) EventType() string     { return "agent.tool.execution_completed" }

// ToolApprovalRequired is raised when a tool requires human approval before execution.
type ToolApprovalRequired struct {
	SessionID  string
	TurnID     string
	ToolName   string
	Args       map[string]any
	occurredAt time.Time
}

func NewToolApprovalRequired(sessionID, turnID, toolName string, args map[string]any) ToolApprovalRequired {
	return ToolApprovalRequired{SessionID: sessionID, TurnID: turnID, ToolName: toolName, Args: args, occurredAt: time.Now()}
}
func (e ToolApprovalRequired) OccurredAt() time.Time { return e.occurredAt }
func (e ToolApprovalRequired) EventType() string     { return "agent.tool.approval_required" }

// SessionCompacted is raised when a session's conversation history is compacted.
type SessionCompacted struct {
	SessionID   string
	TurnsBefore int
	TurnsAfter  int
	occurredAt  time.Time
}

func NewSessionCompacted(sessionID string, turnsBefore, turnsAfter int) SessionCompacted {
	return SessionCompacted{SessionID: sessionID, TurnsBefore: turnsBefore, TurnsAfter: turnsAfter, occurredAt: time.Now()}
}
func (e SessionCompacted) OccurredAt() time.Time { return e.occurredAt }
func (e SessionCompacted) EventType() string     { return "agent.session.compacted" }

// ModelFallbackTriggered is raised when the active model is replaced by a fallback.
type ModelFallbackTriggered struct {
	SessionID  string
	FromModel  string
	ToModel    string
	Reason     string
	occurredAt time.Time
}

func NewModelFallbackTriggered(sessionID, fromModel, toModel, reason string) ModelFallbackTriggered {
	return ModelFallbackTriggered{SessionID: sessionID, FromModel: fromModel, ToModel: toModel, Reason: reason, occurredAt: time.Now()}
}
func (e ModelFallbackTriggered) OccurredAt() time.Time { return e.occurredAt }
func (e ModelFallbackTriggered) EventType() string     { return "agent.model.fallback_triggered" }
