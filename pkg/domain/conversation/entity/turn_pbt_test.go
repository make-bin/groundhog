package entity_test

// Feature: agent-core-enhancements
// Property 6: Turn 记录工具调用顺序
// Validates: Requirements 3.1, 3.2

import (
	"testing"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"pgregory.net/rapid"
)

// TestTurnToolCallOrder verifies that ToolCalls() returns all added tool calls
// in the exact order they were added, regardless of count.
func TestTurnToolCallOrder(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random number of tool calls (0–20)
		n := rapid.IntRange(0, 20).Draw(t, "n")

		turn := entity.NewTurn("turn-id", "user input")

		// Generate and add N tool calls, recording expected names in order
		expectedNames := make([]string, n)
		for i := 0; i < n; i++ {
			name := rapid.StringMatching(`[a-z_]{1,10}`).Draw(t, "tool_name")
			expectedNames[i] = name
			tc := entity.NewToolCall(name, map[string]any{})
			turn.AddToolCall(*tc)
		}

		calls := turn.ToolCalls()

		// Length must equal N
		if len(calls) != n {
			t.Fatalf("expected %d tool calls, got %d", n, len(calls))
		}

		// Order must be preserved
		for i, tc := range calls {
			if tc.ToolName() != expectedNames[i] {
				t.Fatalf("at index %d: expected tool name %q, got %q", i, expectedNames[i], tc.ToolName())
			}
		}
	})
}
