package hook

import (
	"context"
	"sync"
)

// HookPoint defines where in the processing pipeline a hook fires.
type HookPoint string

const (
	HookBeforeMessageReceive HookPoint = "before_message_receive"
	HookAfterMessageReceive  HookPoint = "after_message_receive"
	HookBeforeAgentExecute   HookPoint = "before_agent_execute"
	HookAfterAgentExecute    HookPoint = "after_agent_execute"
	HookBeforeToolExecute    HookPoint = "before_tool_execute"
	HookAfterToolExecute     HookPoint = "after_tool_execute"
)

// HookHandler is a function that handles a hook event.
// data contains context-specific information about the event.
type HookHandler func(ctx context.Context, data interface{}) error

// HookConfig defines a hook loaded from configuration.
type HookConfig struct {
	Name  string    `yaml:"name"`
	Point HookPoint `yaml:"point"`
	// Handler name to look up in a registry (for config-driven hooks)
	Handler string `yaml:"handler"`
}

// HookRegistry manages hook registrations and execution.
type HookRegistry struct {
	mu    sync.RWMutex
	hooks map[HookPoint][]HookHandler
}

// NewHookRegistry creates a new HookRegistry.
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[HookPoint][]HookHandler),
	}
}

// Register adds a handler for the given hook point.
func (r *HookRegistry) Register(point HookPoint, handler HookHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hooks[point] = append(r.hooks[point], handler)
}

// Execute runs all handlers registered for the given hook point.
// Execution stops and returns the first error encountered.
func (r *HookRegistry) Execute(ctx context.Context, point HookPoint, data interface{}) error {
	r.mu.RLock()
	handlers := make([]HookHandler, len(r.hooks[point]))
	copy(handlers, r.hooks[point])
	r.mu.RUnlock()

	for _, h := range handlers {
		if err := h(ctx, data); err != nil {
			return err
		}
	}
	return nil
}

// LoadFromConfig registers named handlers from a config slice.
// namedHandlers maps handler names to HookHandler functions.
func (r *HookRegistry) LoadFromConfig(configs []HookConfig, namedHandlers map[string]HookHandler) {
	for _, cfg := range configs {
		if h, ok := namedHandlers[cfg.Handler]; ok {
			r.Register(cfg.Point, h)
		}
	}
}
