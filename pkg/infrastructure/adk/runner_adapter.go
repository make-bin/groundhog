package adk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/skill"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// MaxToolIterations is the maximum number of tool call iterations in a single turn.
const MaxToolIterations = 10

// ToolEventCallback is called when a tool starts or completes execution.
type ToolEventCallback func(toolCallID, toolName string, args map[string]any, result string, isError bool, durationMs int64, started bool)

// RunnerAdapter orchestrates an agent turn using the ADK layer.
type RunnerAdapter struct {
	ModelAdapter   *ModelAdapter   `inject:""`
	ToolAdapter    *ToolAdapter    `inject:""`
	SessionAdapter *SessionAdapter `inject:""`
	MemorySvc      MemoryService   `inject:""`
	SkillRegistry  *skill.Registry // optional, set before first use
	Logger         logger.Logger   `inject:"logger"`

	toolEventCb ToolEventCallback // set per-turn by StreamTurn
}

// SetToolEventCallback sets the callback for tool execution events (start/done).
// This is called by StreamTurn before each Run() to wire up SSE notifications.
func (r *RunnerAdapter) SetToolEventCallback(cb ToolEventCallback) {
	r.toolEventCb = cb
}

// NewRunnerAdapter creates a new RunnerAdapter.
func NewRunnerAdapter() *RunnerAdapter {
	return &RunnerAdapter{}
}

// Run executes a single agent turn and streams Turn updates via a channel.
// The channel is closed when the turn completes or an error occurs.
func (r *RunnerAdapter) Run(ctx context.Context, sess *agent_session.AgentSession, userInput string) (<-chan *entity.Turn, error) {
	llm, err := r.ModelAdapter.ToADKModel(ctx, sess.ActiveModel())
	if err != nil {
		return nil, fmt.Errorf("runner_adapter: get model: %w", err)
	}

	ch := make(chan *entity.Turn, 1)

	go func() {
		defer close(ch)

		turnID := fmt.Sprintf("turn-%d", time.Now().UnixNano())
		turn := entity.NewTurn(turnID, userInput)
		mc := sess.ActiveModel()
		turn.SetModelUsed(fmt.Sprintf("%s/%s", mc.Provider(), mc.ModelName()))

		// Inject userID into context for memory tools
		if userID := sess.UserID(); userID != "" {
			ctx = WithMemoryUserID(ctx, userID)
		}

		// Check if LLM supports tool calling
		llmWithTools, supportsTools := llm.(LLMWithTools)

		// Build tool schemas (only if LLM supports tools)
		toolSchemas := []ToolSchema{}
		if supportsTools {
			if len(sess.Tools()) > 0 {
				toolSchemas = r.ToolAdapter.ToToolSchemas(sess.Tools())
			}
			// Always append MCP tool schemas
			toolSchemas = append(toolSchemas, r.ToolAdapter.MCPSchemas()...)
			// Append memory tool schemas when memory service is available
			if r.MemorySvc != nil {
				toolSchemas = append(toolSchemas, MemoryToolSchemas()...)
			}
		}

		// Build initial messages
		messages := r.buildMessages(sess, userInput)

		if !supportsTools {
			// Fallback: use legacy string prompt mode
			prompt := r.buildPromptFromMessages(messages)
			streamCh, err := llm.GenerateContentStream(ctx, prompt)
			if err != nil {
				turn.SetResponse(fmt.Sprintf("error: %v", err))
				turn.Complete(vo.NewTokenUsage(0, 0))
				ch <- turn
				return
			}
			var fullResponse string
			for chunk := range streamCh {
				fullResponse += chunk
				partial := entity.NewTurn(turnID, userInput)
				partial.SetResponse(fullResponse)
				ch <- partial
			}
			turn.SetResponse(fullResponse)
			turn.Complete(vo.NewTokenUsage(0, 0))
			ch <- turn
			return
		}

		// Agentic loop
		for i := 0; i < MaxToolIterations; i++ {
			streamCh, _, err := llmWithTools.ChatWithToolsStream(ctx, messages, toolSchemas)
			if err != nil {
				turn.SetResponse(fmt.Sprintf("error: %v", err))
				turn.Complete(vo.NewTokenUsage(0, 0))
				ch <- turn
				return
			}

			// Collect stream: detect tool calls via sentinel or accumulate text
			var textBuf strings.Builder
			var toolCallsResp []LLMToolCall
			for token := range streamCh {
				if strings.HasPrefix(token, "\x00tc:") {
					// Tool calls sentinel
					_ = json.Unmarshal([]byte(token[4:]), &toolCallsResp)
				} else {
					textBuf.WriteString(token)
					// Emit incremental chunk to outer channel
					partial := entity.NewTurn(turnID, userInput)
					partial.SetResponse(textBuf.String())
					ch <- partial
				}
			}

			if len(toolCallsResp) == 0 {
				// Pure text response — done
				turn.SetResponse(textBuf.String())
				turn.Complete(vo.NewTokenUsage(0, 0))
				ch <- turn
				return
			}

			// Build assistant message with tool calls
			assistantMsg := LLMMessage{
				Role:      "assistant",
				ToolCalls: toolCallsResp,
			}
			messages = append(messages, assistantMsg)

			// Execute each tool call
			for _, tc := range toolCallsResp {
				startTime := time.Now()
				toolName := tc.Name

				// Parse arguments
				var args map[string]any
				if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
					args = map[string]any{}
					toolResult := fmt.Sprintf("error: invalid JSON arguments: %v", err)
					messages = append(messages, LLMMessage{
						Role:       "tool",
						Content:    toolResult,
						ToolCallID: tc.ID,
					})
					// Record in turn
					domainTC := entity.NewToolCall(toolName, args)
					domainTC.SetResult(vo.NewToolError(toolResult), time.Since(startTime))
					turn.AddToolCall(*domainTC)
					if r.Logger != nil {
						r.Logger.Info("tool call failed", "tool", toolName, "duration_ms", time.Since(startTime).Milliseconds(), "success", false)
					}
					continue
				}

				// Check BeforeToolCallback
				if r.ToolAdapter.beforeCallback != nil {
					if cbErr := r.ToolAdapter.beforeCallback(ctx, toolName, args); cbErr != nil {
						toolResult := fmt.Sprintf("tool execution denied: %v", cbErr)
						messages = append(messages, LLMMessage{
							Role:       "tool",
							Content:    toolResult,
							ToolCallID: tc.ID,
						})
						domainTC := entity.NewToolCall(toolName, args)
						domainTC.SetResult(vo.NewToolError(toolResult), time.Since(startTime))
						turn.AddToolCall(*domainTC)
						if r.Logger != nil {
							r.Logger.Info("tool call denied", "tool", toolName, "duration_ms", time.Since(startTime).Milliseconds(), "success", false)
						}
						continue
					}
				}

				// Notify tool_start
				if r.toolEventCb != nil {
					r.toolEventCb(tc.ID, toolName, args, "", false, 0, true)
				}

				// Execute tool
				tools := r.ToolAdapter.ToADKTools(sess.Tools())
				if r.MemorySvc != nil {
					saveTool, searchTool := NewMemoryTools(r.MemorySvc)
					tools = append(tools, saveTool, searchTool)
				}
				var toolOutput string
				var toolErr error
				for _, t := range tools {
					if t.Name() == toolName {
						toolOutput, toolErr = t.Execute(ctx, args)
						break
					}
				}

				duration := time.Since(startTime)
				var toolResult string
				var domainResult vo.ToolResult
				if toolErr != nil {
					toolResult = fmt.Sprintf("error: %v", toolErr)
					domainResult = vo.NewToolError(toolErr.Error())
					if r.Logger != nil {
						r.Logger.Info("tool call completed", "tool", toolName, "duration_ms", duration.Milliseconds(), "success", false)
					}
				} else {
					if toolOutput == "" && toolErr == nil {
						// Tool not found
						toolResult = fmt.Sprintf("tool not found: %s", toolName)
						domainResult = vo.NewToolError(toolResult)
					} else {
						toolResult = toolOutput
						domainResult = vo.NewToolResult(toolOutput)
					}
					if r.Logger != nil {
						r.Logger.Info("tool call completed", "tool", toolName, "duration_ms", duration.Milliseconds(), "success", true)
					}
				}

				// Notify tool_done
				if r.toolEventCb != nil {
					r.toolEventCb(tc.ID, toolName, args, toolResult, toolErr != nil || strings.HasPrefix(toolResult, "tool not found"), duration.Milliseconds(), false)
				}

				messages = append(messages, LLMMessage{
					Role:       "tool",
					Content:    toolResult,
					ToolCallID: tc.ID,
				})

				domainTC := entity.NewToolCall(toolName, args)
				domainTC.SetResult(domainResult, duration)
				turn.AddToolCall(*domainTC)
			}
		}

		// MaxToolIterations reached — do one final streaming call to get text response
		finalCh, _, err := llmWithTools.ChatWithToolsStream(ctx, messages, toolSchemas)
		if err != nil {
			turn.SetResponse(fmt.Sprintf("error: %v", err))
		} else {
			var textBuf strings.Builder
			for token := range finalCh {
				if !strings.HasPrefix(token, "\x00tc:") {
					textBuf.WriteString(token)
					partial := entity.NewTurn(turnID, userInput)
					partial.SetResponse(textBuf.String())
					ch <- partial
				}
			}
			turn.SetResponse(textBuf.String())
		}
		turn.Complete(vo.NewTokenUsage(0, 0))
		ch <- turn
	}()

	return ch, nil
}

// buildMessages constructs the structured message list for the LLM.
// isSummary turns are formatted with a [Context Summary] prefix.
func (r *RunnerAdapter) buildMessages(sess *agent_session.AgentSession, userInput string) []LLMMessage {
	var messages []LLMMessage

	// 1. Skills system prompt
	if r.SkillRegistry != nil {
		skillPrompt := r.SkillRegistry.ResolvePrompt(sess.Skills())
		if skillPrompt != "" {
			messages = append(messages, LLMMessage{Role: "system", Content: skillPrompt})
		}
	}

	// 2. Session-level system prompt
	if sp := sess.SystemPrompt().Render(); sp != "" {
		messages = append(messages, LLMMessage{Role: "system", Content: sp})
	}

	// 3. Conversation history (last 10 turns)
	turns := sess.Turns()
	start := 0
	if len(turns) > 10 {
		start = len(turns) - 10
	}
	for _, t := range turns[start:] {
		if t.IsSummary() {
			messages = append(messages, LLMMessage{
				Role:    "assistant",
				Content: "[Context Summary]\n" + t.Response(),
			})
		} else {
			if t.UserInput() != "" {
				messages = append(messages, LLMMessage{Role: "user", Content: t.UserInput()})
			}
			if t.Response() != "" {
				messages = append(messages, LLMMessage{Role: "assistant", Content: t.Response()})
			}
		}
	}

	// 4. Current user input
	messages = append(messages, LLMMessage{Role: "user", Content: userInput})

	return messages
}

// buildPromptFromMessages converts a message list to a plain text prompt string (legacy fallback).
func (r *RunnerAdapter) buildPromptFromMessages(messages []LLMMessage) string {
	var parts []string
	for _, m := range messages {
		switch m.Role {
		case "system":
			parts = append(parts, m.Content)
		case "user":
			parts = append(parts, fmt.Sprintf("User: %s", m.Content))
		case "assistant":
			parts = append(parts, fmt.Sprintf("Assistant: %s", m.Content))
		}
	}
	// Append the trailing "Assistant:" prompt
	parts = append(parts, "Assistant:")
	return strings.Join(parts, "\n\n")
}
