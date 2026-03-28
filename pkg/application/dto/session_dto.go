package dto

import "time"

// CreateSessionRequest is the request DTO for creating a new session.
type CreateSessionRequest struct {
	AgentID      string   `json:"agent_id" binding:"required"`
	UserID       string   `json:"user_id" binding:"required"`
	Provider     string   `json:"provider" binding:"required"`
	ModelName    string   `json:"model_name" binding:"required"`
	Temperature  float64  `json:"temperature"`
	MaxTokens    int      `json:"max_tokens"`
	Skills       []string `json:"skills"`
	SystemPrompt string   `json:"system_prompt"`
}

// SendMessageRequest is the request DTO for sending a message to a session.
type SendMessageRequest struct {
	UserInput string `json:"user_input" binding:"required"`
}

// SessionListRequest is the request DTO for listing sessions.
type SessionListRequest struct {
	UserID  string `json:"user_id" form:"user_id"`
	AgentID string `json:"agent_id" form:"agent_id"`
	State   *int   `json:"state" form:"state"`
	Offset  int    `json:"offset" form:"offset"`
	Limit   int    `json:"limit" form:"limit"`
}

// TurnResponse is the response DTO for a single conversation turn.
type TurnResponse struct {
	ID          string    `json:"id"`
	UserInput   string    `json:"user_input"`
	Response    string    `json:"response"`
	ModelUsed   string    `json:"model_used"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

// SessionResponse is the response DTO for a session.
type SessionResponse struct {
	ID           string         `json:"id"`
	AgentID      string         `json:"agent_id"`
	UserID       string         `json:"user_id"`
	State        string         `json:"state"`
	ActiveModel  string         `json:"active_model"`
	Turns        []TurnResponse `json:"turns"`
	CreatedAt    time.Time      `json:"created_at"`
	LastActiveAt time.Time      `json:"last_active_at"`
}

// SessionListResponse is the response DTO for a list of sessions.
type SessionListResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
	Total    int                `json:"total"`
	Offset   int                `json:"offset"`
	Limit    int                `json:"limit"`
}

// ApprovalRequiredEvent is pushed over SSE when a tool needs user confirmation.
type ApprovalRequiredEvent struct {
	ApprovalID string         `json:"approval_id"`
	SessionID  string         `json:"session_id"`
	ToolName   string         `json:"tool_name"`
	Args       map[string]any `json:"args"`
}

// ApprovalRequest is the request DTO for resolving a pending approval.
type ApprovalRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approve deny"`
}

// ToolCallEvent carries tool execution lifecycle information over SSE.
type ToolCallEvent struct {
	ToolCallID string         `json:"tool_call_id"`
	ToolName   string         `json:"tool_name"`
	Args       map[string]any `json:"args,omitempty"`
	Result     string         `json:"result,omitempty"`
	IsError    bool           `json:"is_error,omitempty"`
	DurationMs int64          `json:"duration_ms,omitempty"`
}

// StreamEvent wraps all possible SSE event types.
type StreamEvent struct {
	Type string `json:"type"` // "chunk" | "tool_start" | "tool_done" | "approval_required" | "done" | "error"
	// chunk
	Chunk string `json:"chunk,omitempty"`
	// tool_start / tool_done
	Tool *ToolCallEvent `json:"tool,omitempty"`
	// approval_required
	Approval *ApprovalRequiredEvent `json:"approval,omitempty"`
	// done
	Turn *TurnResponse `json:"turn,omitempty"`
	// error
	Error string `json:"error,omitempty"`
}

// PendingApprovalResponse is the response DTO for a pending tool approval.
type PendingApprovalResponse struct {
	ApprovalID string         `json:"approval_id"`
	SessionID  string         `json:"session_id"`
	ToolName   string         `json:"tool_name"`
	Args       map[string]any `json:"args"`
	CreatedAt  time.Time      `json:"created_at"`
}
