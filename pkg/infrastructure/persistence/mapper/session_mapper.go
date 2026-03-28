// @AI_GENERATED
package mapper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// modelConfigJSON is the JSON representation of vo.ModelConfig.
type modelConfigJSON struct {
	Provider      int      `json:"provider"`
	ModelName     string   `json:"model_name"`
	Temperature   float64  `json:"temperature"`
	MaxTokens     int      `json:"max_tokens"`
	FallbackChain []string `json:"fallback_chain"`
	AuthProfile   string   `json:"auth_profile"`
}

// tokenUsageJSON is the JSON representation of vo.TokenUsage.
type tokenUsageJSON struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

// toolCallJSON is the JSON representation of entity.ToolCall.
type toolCallJSON struct {
	ToolName string         `json:"tool_name"`
	Args     map[string]any `json:"args"`
	Output   string         `json:"output"`
	IsError  bool           `json:"is_error"`
	ErrMsg   string         `json:"err_msg"`
	Duration int64          `json:"duration_ns"`
	Approved bool           `json:"approved"`
}

// DomainToSessionPO converts an AgentSession aggregate to a SessionPO.
func DomainToSessionPO(session *agent_session.AgentSession) (*po.SessionPO, error) {
	activeModelJSON, err := marshalModelConfig(session.ActiveModel())
	if err != nil {
		return nil, fmt.Errorf("marshal active_model: %w", err)
	}

	metadataJSON, err := json.Marshal(session.Metadata())
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	turns, err := domainTurnsToTurnPOs(session.Turns())
	if err != nil {
		return nil, err
	}

	return &po.SessionPO{
		SessionID:   session.ID().Value(),
		AgentID:     session.AgentID().Value(),
		UserID:      session.UserID(),
		ActiveModel: string(activeModelJSON),
		State:       int(session.State()),
		Metadata:    string(metadataJSON),
		Turns:       turns,
	}, nil
}

// SessionPOToDomain converts a SessionPO to an AgentSession aggregate.
func SessionPOToDomain(sessionPO *po.SessionPO) (*agent_session.AgentSession, error) {
	sessionID, err := vo.NewSessionID(sessionPO.SessionID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct session_id: %w", err)
	}

	agentID, err := vo.NewAgentID(sessionPO.AgentID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct agent_id: %w", err)
	}

	activeModel, err := unmarshalModelConfig(sessionPO.ActiveModel)
	if err != nil {
		return nil, fmt.Errorf("unmarshal active_model: %w", err)
	}

	var metadata map[string]any
	if sessionPO.Metadata != "" {
		if err := json.Unmarshal([]byte(sessionPO.Metadata), &metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}
	if metadata == nil {
		metadata = map[string]any{}
	}

	turns, err := turnPOsToDomainTurns(sessionPO.Turns)
	if err != nil {
		return nil, err
	}

	return agent_session.ReconstructAgentSession(
		sessionID,
		agentID,
		sessionPO.UserID,
		turns,
		activeModel,
		[]entity.ToolDefinition{},
		vo.NewPrompt("", nil),
		[]string{},
		vo.SessionState(sessionPO.State),
		sessionPO.CreatedAt,
		sessionPO.UpdatedAt,
		metadata,
		0,
	), nil
}

// marshalModelConfig serializes a vo.ModelConfig to JSON bytes.
func marshalModelConfig(cfg vo.ModelConfig) ([]byte, error) {
	return json.Marshal(modelConfigJSON{
		Provider:      int(cfg.Provider()),
		ModelName:     cfg.ModelName(),
		Temperature:   cfg.Temperature(),
		MaxTokens:     cfg.MaxTokens(),
		FallbackChain: cfg.FallbackChain(),
		AuthProfile:   cfg.AuthProfile(),
	})
}

// unmarshalModelConfig deserializes a JSON string to vo.ModelConfig.
func unmarshalModelConfig(raw string) (vo.ModelConfig, error) {
	if raw == "" {
		return vo.ModelConfig{}, nil
	}
	var j modelConfigJSON
	if err := json.Unmarshal([]byte(raw), &j); err != nil {
		return vo.ModelConfig{}, err
	}
	return vo.NewModelConfig(
		vo.ProviderType(j.Provider),
		j.ModelName,
		j.Temperature,
		j.MaxTokens,
		j.FallbackChain,
		j.AuthProfile,
	)
}

// domainTurnsToTurnPOs converts a slice of entity.Turn to []po.TurnPO.
func domainTurnsToTurnPOs(turns []entity.Turn) ([]po.TurnPO, error) {
	result := make([]po.TurnPO, 0, len(turns))
	for i := range turns {
		t := &turns[i]
		tokenUsageBytes, err := json.Marshal(tokenUsageJSON{
			PromptTokens:     t.TokenUsage().PromptTokens(),
			CompletionTokens: t.TokenUsage().CompletionTokens(),
		})
		if err != nil {
			return nil, fmt.Errorf("marshal token_usage for turn %s: %w", t.ID(), err)
		}

		toolCallsBytes, err := marshalToolCalls(t.ToolCalls())
		if err != nil {
			return nil, fmt.Errorf("marshal tool_calls for turn %s: %w", t.ID(), err)
		}

		result = append(result, po.TurnPO{
			TurnID:      t.ID(),
			UserInput:   t.UserInput(),
			Response:    t.Response(),
			ModelUsed:   t.ModelUsed(),
			TokenUsage:  string(tokenUsageBytes),
			ToolCalls:   string(toolCallsBytes),
			StartedAt:   t.StartedAt(),
			CompletedAt: t.CompletedAt(),
		})
	}
	return result, nil
}

// turnPOsToDomainTurns converts a slice of po.TurnPO to []entity.Turn.
func turnPOsToDomainTurns(turnPOs []po.TurnPO) ([]entity.Turn, error) {
	result := make([]entity.Turn, 0, len(turnPOs))
	for i := range turnPOs {
		t := &turnPOs[i]

		var tokenUsage tokenUsageJSON
		if t.TokenUsage != "" {
			if err := json.Unmarshal([]byte(t.TokenUsage), &tokenUsage); err != nil {
				return nil, fmt.Errorf("unmarshal token_usage for turn %s: %w", t.TurnID, err)
			}
		}

		toolCalls, err := unmarshalToolCalls(t.ToolCalls)
		if err != nil {
			return nil, fmt.Errorf("unmarshal tool_calls for turn %s: %w", t.TurnID, err)
		}

		turn := entity.ReconstructTurn(
			t.TurnID,
			t.UserInput,
			t.Response,
			t.ModelUsed,
			toolCalls,
			vo.NewTokenUsage(tokenUsage.PromptTokens, tokenUsage.CompletionTokens),
			t.StartedAt,
			t.CompletedAt,
			false, // isSummary: existing persisted turns are not summaries
		)
		result = append(result, *turn)
	}
	return result, nil
}

// marshalToolCalls serializes a slice of entity.ToolCall to JSON bytes.
func marshalToolCalls(calls []entity.ToolCall) ([]byte, error) {
	jCalls := make([]toolCallJSON, 0, len(calls))
	for i := range calls {
		c := &calls[i]
		result := c.Result()
		jCalls = append(jCalls, toolCallJSON{
			ToolName: c.ToolName(),
			Args:     c.Args(),
			Output:   result.Output(),
			IsError:  result.IsError(),
			ErrMsg:   result.ErrMsg(),
			Duration: c.Duration().Nanoseconds(),
			Approved: c.Approved(),
		})
	}
	return json.Marshal(jCalls)
}

// unmarshalToolCalls deserializes a JSON string to []entity.ToolCall.
func unmarshalToolCalls(raw string) ([]entity.ToolCall, error) {
	if raw == "" {
		return []entity.ToolCall{}, nil
	}
	var jCalls []toolCallJSON
	if err := json.Unmarshal([]byte(raw), &jCalls); err != nil {
		return nil, err
	}

	calls := make([]entity.ToolCall, 0, len(jCalls))
	for _, j := range jCalls {
		tc := entity.NewToolCall(j.ToolName, j.Args)
		var result vo.ToolResult
		if j.IsError {
			result = vo.NewToolError(j.ErrMsg)
		} else {
			result = vo.NewToolResult(j.Output)
		}
		tc.SetResult(result, time.Duration(j.Duration))
		if j.Approved {
			tc.Approve()
		}
		calls = append(calls, *tc)
	}
	return calls, nil
}

// @AI_GENERATED: end
