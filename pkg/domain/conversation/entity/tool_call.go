// @AI_GENERATED
package entity

import (
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
)

// ToolCall is an entity representing a single tool invocation during a turn.
type ToolCall struct {
	toolName string
	args     map[string]any
	result   vo.ToolResult
	duration time.Duration
	approved bool
}

// NewToolCall creates a new ToolCall with the given tool name and arguments.
func NewToolCall(toolName string, args map[string]any) *ToolCall {
	a := make(map[string]any, len(args))
	for k, v := range args {
		a[k] = v
	}
	return &ToolCall{
		toolName: toolName,
		args:     a,
	}
}

// ToolName returns the name of the tool that was called.
func (tc *ToolCall) ToolName() string { return tc.toolName }

// Args returns a copy of the tool call arguments.
func (tc *ToolCall) Args() map[string]any {
	a := make(map[string]any, len(tc.args))
	for k, v := range tc.args {
		a[k] = v
	}
	return a
}

// Result returns the tool execution result.
func (tc *ToolCall) Result() vo.ToolResult { return tc.result }

// Duration returns how long the tool execution took.
func (tc *ToolCall) Duration() time.Duration { return tc.duration }

// Approved returns whether the tool call was approved.
func (tc *ToolCall) Approved() bool { return tc.approved }

// SetResult records the tool execution result and duration.
func (tc *ToolCall) SetResult(result vo.ToolResult, duration time.Duration) {
	tc.result = result
	tc.duration = duration
}

// Approve marks the tool call as approved.
func (tc *ToolCall) Approve() {
	tc.approved = true
}

// @AI_GENERATED: end
