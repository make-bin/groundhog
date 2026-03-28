package vo

import "fmt"

// PayloadKind identifies the type of cron payload.
type PayloadKind string

const (
	PayloadKindSystemEvent PayloadKind = "systemEvent"
	PayloadKindAgentTurn   PayloadKind = "agentTurn"
)

// CronPayload holds the execution content for a cron job.
// It is immutable after construction.
type CronPayload struct {
	kind           PayloadKind
	text           string // systemEvent
	message        string // agentTurn
	model          string // agentTurn optional
	thinking       bool   // agentTurn optional
	timeoutSeconds int    // agentTurn optional
	lightContext   bool   // agentTurn optional
}

// NewCronPayloadSystemEvent creates a systemEvent payload.
// text must be non-empty.
func NewCronPayloadSystemEvent(text string) (CronPayload, error) {
	if text == "" {
		return CronPayload{}, fmt.Errorf("systemEvent payload text must not be empty")
	}
	return CronPayload{kind: PayloadKindSystemEvent, text: text}, nil
}

// NewCronPayloadAgentTurn creates an agentTurn payload.
// message must be non-empty; timeoutSeconds must be >= 0.
func NewCronPayloadAgentTurn(message, model string, thinking bool, timeoutSeconds int, lightContext bool) (CronPayload, error) {
	if message == "" {
		return CronPayload{}, fmt.Errorf("agentTurn payload message must not be empty")
	}
	if timeoutSeconds < 0 {
		return CronPayload{}, fmt.Errorf("agentTurn payload timeoutSeconds must be >= 0, got %d", timeoutSeconds)
	}
	return CronPayload{
		kind:           PayloadKindAgentTurn,
		message:        message,
		model:          model,
		thinking:       thinking,
		timeoutSeconds: timeoutSeconds,
		lightContext:   lightContext,
	}, nil
}

// Getters

func (p CronPayload) Kind() PayloadKind   { return p.kind }
func (p CronPayload) Text() string        { return p.text }
func (p CronPayload) Message() string     { return p.message }
func (p CronPayload) Model() string       { return p.model }
func (p CronPayload) Thinking() bool      { return p.thinking }
func (p CronPayload) TimeoutSeconds() int { return p.timeoutSeconds }
func (p CronPayload) LightContext() bool  { return p.lightContext }
