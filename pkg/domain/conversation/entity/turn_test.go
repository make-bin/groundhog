package entity_test

import (
	"testing"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/stretchr/testify/assert"
)

func TestNewTurn_IsSummaryFalse(t *testing.T) {
	turn := entity.NewTurn("id-1", "hello")
	assert.False(t, turn.IsSummary())
	assert.Equal(t, "hello", turn.UserInput())
}

func TestNewSummaryTurn_IsSummaryTrue(t *testing.T) {
	turn := entity.NewSummaryTurn("sum-1", "This is a summary.")
	assert.True(t, turn.IsSummary())
	assert.Equal(t, "", turn.UserInput())
	assert.Equal(t, "This is a summary.", turn.Response())
	assert.Equal(t, "sum-1", turn.ID())
}

func TestAddToolCall_PreservesOrder(t *testing.T) {
	turn := entity.NewTurn("id-2", "input")

	tc1 := entity.NewToolCall("tool_a", map[string]any{"x": 1})
	tc2 := entity.NewToolCall("tool_b", map[string]any{"y": 2})
	tc3 := entity.NewToolCall("tool_c", map[string]any{"z": 3})

	turn.AddToolCall(*tc1)
	turn.AddToolCall(*tc2)
	turn.AddToolCall(*tc3)

	calls := turn.ToolCalls()
	assert.Len(t, calls, 3)
	assert.Equal(t, "tool_a", calls[0].ToolName())
	assert.Equal(t, "tool_b", calls[1].ToolName())
	assert.Equal(t, "tool_c", calls[2].ToolName())
}

func TestToolCalls_ReturnsCopy(t *testing.T) {
	turn := entity.NewTurn("id-3", "input")
	tc := entity.NewToolCall("tool_a", map[string]any{})
	turn.AddToolCall(*tc)

	calls := turn.ToolCalls()
	// Mutating the returned slice should not affect the turn's internal state.
	calls[0] = *entity.NewToolCall("mutated_tool", map[string]any{})

	assert.Equal(t, "tool_a", turn.ToolCalls()[0].ToolName())
}

func TestReconstructTurn_IsSummaryField(t *testing.T) {
	now := time.Now()
	usage := vo.NewTokenUsage(10, 20)

	// Reconstruct a summary turn
	summaryTurn := entity.ReconstructTurn(
		"s-1", "", "summary text", "gpt-4",
		nil, usage, now, now,
		true,
	)
	assert.True(t, summaryTurn.IsSummary())
	assert.Equal(t, "summary text", summaryTurn.Response())

	// Reconstruct a normal turn
	normalTurn := entity.ReconstructTurn(
		"n-1", "user input", "response", "gpt-4",
		nil, usage, now, now,
		false,
	)
	assert.False(t, normalTurn.IsSummary())
}
