package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/make-bin/groundhog/pkg/application/approval"
	"github.com/make-bin/groundhog/pkg/application/assembler"
	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/application/eventbus"
	"github.com/make-bin/groundhog/pkg/application/hook"
	"github.com/make-bin/groundhog/pkg/domain/conversation/aggregate/agent_session"
	"github.com/make-bin/groundhog/pkg/domain/conversation/entity"
	"github.com/make-bin/groundhog/pkg/domain/conversation/event"
	"github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	conversation_service "github.com/make-bin/groundhog/pkg/domain/conversation/service"
	"github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/adk"
	"github.com/make-bin/groundhog/pkg/utils/config"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// AgentAppService defines the application service interface for agent session management.
type AgentAppService interface {
	CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.SessionResponse, error)
	ExecuteTurn(ctx context.Context, sessionID vo.SessionID, userInput string) (*dto.TurnResponse, error)
	StreamTurn(ctx context.Context, sessionID vo.SessionID, userInput string) <-chan *dto.StreamEvent
	ResolveApproval(approvalID string, decision string) error
	ListApprovals(sessionID string) []*dto.PendingApprovalResponse
	ExecuteWorkflow(ctx context.Context, sessionIDs []vo.SessionID, workflowType string, initialInput string) (string, error)
	GetSession(ctx context.Context, id vo.SessionID) (*dto.SessionResponse, error)
	ListSessions(ctx context.Context, filter dto.SessionListRequest) (*dto.SessionListResponse, error)
	DeleteSession(ctx context.Context, id vo.SessionID) error
}

// memoryRecallPrompt is injected into the session system prompt when memory is enabled.
const memoryRecallPrompt = `## Memory Recall
Before responding, use memory_search to retrieve relevant memories about the user.
After important interactions, use memory_save to store key information.`

type agentAppService struct {
	SessionRepo    repository.SessionRepository           `inject:""`
	RunnerAdapter  *adk.RunnerAdapter                     `inject:""`
	WorkflowRunner *adk.WorkflowRunner                    `inject:""`
	EventBus       eventbus.EventBus                      `inject:""`
	HookRegistry   *hook.HookRegistry                     `inject:""`
	ApprovalMgr    *approval.Manager                      `inject:""`
	CompactionSvc  conversation_service.CompactionService `inject:""`
	Logger         logger.Logger                          `inject:"logger"`
	cfg            *config.AppConfig
}

// NewAgentAppService creates a new AgentAppService. Dependencies are injected via struct tags.
func NewAgentAppService(cfg *config.AppConfig) AgentAppService {
	return &agentAppService{cfg: cfg}
}

func (s *agentAppService) CreateSession(ctx context.Context, req *dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	sessionID, err := vo.NewSessionID(fmt.Sprintf("sess-%d", time.Now().UnixNano()))
	if err != nil {
		return nil, err
	}
	agentID, err := vo.NewAgentID(req.AgentID)
	if err != nil {
		return nil, err
	}
	provider := vo.ProviderGemini
	switch req.Provider {
	case "openai":
		provider = vo.ProviderOpenAI
	case "anthropic":
		provider = vo.ProviderAnthropic
	case "ollama":
		provider = vo.ProviderOllama
	case "groq":
		provider = vo.ProviderGroq
	case "openai_compat":
		provider = vo.ProviderOpenAICompat
	}
	modelCfg, err := vo.NewModelConfig(provider, req.ModelName, req.Temperature, req.MaxTokens, nil, "")
	if err != nil {
		return nil, err
	}
	systemPrompt := req.SystemPrompt
	if s.cfg != nil && s.cfg.Memory.Enabled && req.UserID != "" {
		systemPrompt += "\n\n" + memoryRecallPrompt
	}
	sess, err := agent_session.NewAgentSession(sessionID, agentID, req.UserID, modelCfg, vo.NewPrompt(systemPrompt, nil))
	if err != nil {
		return nil, err
	}
	if len(req.Skills) > 0 {
		sess.SetSkills(req.Skills)
	}
	if err := s.SessionRepo.Create(ctx, sess); err != nil {
		return nil, err
	}
	return assembler.ToSessionResponse(sess), nil
}

func (s *agentAppService) ExecuteTurn(ctx context.Context, sessionID vo.SessionID, userInput string) (*dto.TurnResponse, error) {
	sess, err := s.SessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if sess.NeedsCompaction(50) && s.CompactionSvc != nil {
		turnsBefore := len(sess.Turns())
		if err := s.CompactionSvc.Compact(ctx, sess, 5); err == nil {
			s.EventBus.Publish(event.NewSessionCompacted(sessionID.Value(), turnsBefore, len(sess.Turns())))
		} else {
			s.Logger.Warn("compaction failed", "error", err, "sessionID", sessionID.Value())
		}
	}

	if s.HookRegistry != nil {
		if err := s.HookRegistry.Execute(ctx, hook.HookBeforeAgentExecute, map[string]interface{}{
			"sessionID": sessionID.Value(),
			"userInput": userInput,
		}); err != nil {
			return nil, fmt.Errorf("agent_app_service: before_agent_execute hook: %w", err)
		}
	}

	s.EventBus.Publish(event.NewAgentTurnStarted(sessionID.Value(), "", userInput))

	turnCh, err := s.RunnerAdapter.Run(ctx, sess, userInput)
	if err != nil {
		return nil, fmt.Errorf("agent_app_service: run: %w", err)
	}

	var lastTurn *entity.Turn
	for turn := range turnCh {
		lastTurn = turn
		s.EventBus.Publish(event.NewAgentTurnCompleted(
			sessionID.Value(),
			turn.ID(),
			turn.Response(),
			turn.TokenUsage().PromptTokens(),
			turn.TokenUsage().CompletionTokens(),
		))
	}

	if lastTurn != nil {
		sess.AddTurn(*lastTurn)
		if err := s.SessionRepo.Update(ctx, sess); err != nil {
			return nil, fmt.Errorf("agent_app_service: update session: %w", err)
		}
	}

	if s.HookRegistry != nil {
		if err := s.HookRegistry.Execute(ctx, hook.HookAfterAgentExecute, map[string]interface{}{
			"sessionID": sessionID.Value(),
			"userInput": userInput,
		}); err != nil {
			return nil, fmt.Errorf("agent_app_service: after_agent_execute hook: %w", err)
		}
	}

	if lastTurn == nil {
		return &dto.TurnResponse{UserInput: userInput}, nil
	}
	return &dto.TurnResponse{
		ID:          lastTurn.ID(),
		UserInput:   lastTurn.UserInput(),
		Response:    lastTurn.Response(),
		ModelUsed:   lastTurn.ModelUsed(),
		StartedAt:   lastTurn.StartedAt(),
		CompletedAt: lastTurn.CompletedAt(),
	}, nil
}

// StreamTurn streams turn events as a single channel of StreamEvent.
// Event types: "chunk" (token delta), "approval_required" (waiting for user), "done" (final turn), "error".
func (s *agentAppService) StreamTurn(ctx context.Context, sessionID vo.SessionID, userInput string) <-chan *dto.StreamEvent {
	eventCh := make(chan *dto.StreamEvent, 32)

	go func() {
		defer close(eventCh)

		send := func(e *dto.StreamEvent) {
			select {
			case eventCh <- e:
			case <-ctx.Done():
			}
		}

		sess, err := s.SessionRepo.FindByID(ctx, sessionID)
		if err != nil {
			send(&dto.StreamEvent{Type: "error", Error: err.Error()})
			return
		}

		if sess.NeedsCompaction(50) && s.CompactionSvc != nil {
			turnsBefore := len(sess.Turns())
			if err := s.CompactionSvc.Compact(ctx, sess, 5); err == nil {
				s.EventBus.Publish(event.NewSessionCompacted(sessionID.Value(), turnsBefore, len(sess.Turns())))
			} else {
				s.Logger.Warn("compaction failed", "error", err, "sessionID", sessionID.Value())
			}
		}

		// Install tool event callback to stream tool_start / tool_done events
		s.RunnerAdapter.SetToolEventCallback(func(toolCallID, toolName string, args map[string]any, result string, isError bool, durationMs int64, started bool) {
			evtType := "tool_done"
			if started {
				evtType = "tool_start"
			}
			send(&dto.StreamEvent{
				Type: evtType,
				Tool: &dto.ToolCallEvent{
					ToolCallID: toolCallID,
					ToolName:   toolName,
					Args:       args,
					Result:     result,
					IsError:    isError,
					DurationMs: durationMs,
				},
			})
		})

		// Install approval gate on the tool adapter before running
		if s.ApprovalMgr != nil {
			s.RunnerAdapter.ToolAdapter.SetBeforeToolCallback(func(tctx context.Context, toolName string, args map[string]any) error {
				if s.requiresApproval(toolName) {
					// Create pending approval and notify SSE clients before blocking
					approvalID := fmt.Sprintf("appr-%d", time.Now().UnixNano())
					pa := &approval.PendingApproval{
						ID:        approvalID,
						SessionID: sessionID.Value(),
						ToolName:  toolName,
						Args:      args,
					}
					// Notify SSE stream about pending approval
					send(&dto.StreamEvent{
						Type: "approval_required",
						Approval: &dto.ApprovalRequiredEvent{
							ApprovalID: approvalID,
							SessionID:  sessionID.Value(),
							ToolName:   toolName,
							Args:       args,
						},
					})
					// Block until user decides
					approved, _, aerr := s.ApprovalMgr.RequestWithID(tctx, pa)
					if aerr != nil {
						return fmt.Errorf("approval cancelled: %w", aerr)
					}
					if !approved {
						return fmt.Errorf("tool %s denied by user", toolName)
					}
				}
				return nil
			})
		}

		turnCh, err := s.RunnerAdapter.Run(ctx, sess, userInput)
		if err != nil {
			send(&dto.StreamEvent{Type: "error", Error: err.Error()})
			return
		}

		var lastTurn *entity.Turn
		var prevLen int
		for turn := range turnCh {
			full := turn.Response()
			if len(full) > prevLen {
				send(&dto.StreamEvent{Type: "chunk", Chunk: full[prevLen:]})
				prevLen = len(full)
			}
			lastTurn = turn
		}

		if lastTurn != nil {
			sess.AddTurn(*lastTurn)
			if err := s.SessionRepo.Update(ctx, sess); err != nil {
				send(&dto.StreamEvent{Type: "error", Error: err.Error()})
				return
			}
			send(&dto.StreamEvent{
				Type: "done",
				Turn: &dto.TurnResponse{
					ID:          lastTurn.ID(),
					UserInput:   lastTurn.UserInput(),
					Response:    lastTurn.Response(),
					ModelUsed:   lastTurn.ModelUsed(),
					StartedAt:   lastTurn.StartedAt(),
					CompletedAt: lastTurn.CompletedAt(),
				},
			})
		}
	}()

	return eventCh
}

// requiresApproval returns true if the given tool name needs human confirmation.
// It checks:
//  1. Built-in dangerous tools (always require approval)
//  2. MCP server-level require_approval flag
//  3. MCP server-level dangerous_tools list
func (s *agentAppService) requiresApproval(toolName string) bool {
	// Built-in dangerous tools
	builtinDangerous := map[string]bool{
		"bash_exec":   true,
		"file_write":  true,
		"file_delete": true,
	}
	if builtinDangerous[toolName] {
		return true
	}

	// Check MCP server configs
	if s.cfg == nil {
		return false
	}
	for _, srv := range s.cfg.MCP.Servers {
		// Server-level: all tools require approval
		if srv.RequireApproval {
			// Check if this tool belongs to this server (bare name or prefixed)
			prefix := srv.Name + "__"
			if strings.HasPrefix(toolName, prefix) || s.isMCPToolOfServer(toolName, srv.Name) {
				return true
			}
		}
		// Tool-level: specific dangerous tools
		for _, dt := range srv.DangerousTools {
			if toolName == dt || toolName == srv.Name+"__"+dt {
				return true
			}
		}
	}
	return false
}

// isMCPToolOfServer checks if a bare tool name was registered from the given server.
// Since MCP tools are registered with both bare and prefixed names, we check the prefix form.
func (s *agentAppService) isMCPToolOfServer(toolName, serverName string) bool {
	// The tool adapter registers both "server__tool" and "tool" forms.
	// We can't easily reverse-map bare names to servers here, so we rely on
	// the dangerous_tools list for bare-name matching.
	_ = serverName
	return false
}

func (s *agentAppService) ResolveApproval(approvalID string, decision string) error {
	if s.ApprovalMgr == nil {
		return fmt.Errorf("approval manager not configured")
	}
	d := approval.DecisionDeny
	if decision == "approve" {
		d = approval.DecisionApprove
	}
	return s.ApprovalMgr.Resolve(approvalID, d)
}

func (s *agentAppService) ListApprovals(sessionID string) []*dto.PendingApprovalResponse {
	if s.ApprovalMgr == nil {
		return nil
	}
	pending := s.ApprovalMgr.ListPending(sessionID)
	result := make([]*dto.PendingApprovalResponse, 0, len(pending))
	for _, pa := range pending {
		result = append(result, &dto.PendingApprovalResponse{
			ApprovalID: pa.ID,
			SessionID:  pa.SessionID,
			ToolName:   pa.ToolName,
			Args:       pa.Args,
			CreatedAt:  pa.CreatedAt,
		})
	}
	return result
}

func (s *agentAppService) ExecuteWorkflow(ctx context.Context, sessionIDs []vo.SessionID, workflowType string, initialInput string) (string, error) {
	sessions := make([]*agent_session.AgentSession, 0, len(sessionIDs))
	for _, id := range sessionIDs {
		sess, err := s.SessionRepo.FindByID(ctx, id)
		if err != nil {
			return "", fmt.Errorf("load session %s: %w", id.Value(), err)
		}
		sessions = append(sessions, sess)
	}

	switch adk.WorkflowType(workflowType) {
	case adk.WorkflowSequential:
		return s.WorkflowRunner.RunSequential(ctx, sessions, initialInput)
	case adk.WorkflowLoop:
		if len(sessions) == 0 {
			return "", fmt.Errorf("loop workflow requires at least one session")
		}
		return s.WorkflowRunner.RunLoop(ctx, sessions[0], initialInput, 3)
	case adk.WorkflowParallel:
		results, err := s.WorkflowRunner.RunParallel(ctx, sessions, initialInput)
		if err != nil {
			return "", err
		}
		combined := ""
		for i, r := range results {
			if i > 0 {
				combined += "\n---\n"
			}
			combined += r
		}
		return combined, nil
	default:
		return "", fmt.Errorf("unknown workflow type: %s", workflowType)
	}
}

func (s *agentAppService) GetSession(ctx context.Context, id vo.SessionID) (*dto.SessionResponse, error) {
	sess, err := s.SessionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return assembler.ToSessionResponse(sess), nil
}

func (s *agentAppService) ListSessions(ctx context.Context, req dto.SessionListRequest) (*dto.SessionListResponse, error) {
	filter := repository.SessionFilter{}
	if req.UserID != "" {
		filter.UserID = &req.UserID
	}
	if req.AgentID != "" {
		agentID, err := vo.NewAgentID(req.AgentID)
		if err != nil {
			return nil, err
		}
		filter.AgentID = &agentID
	}
	if req.State != nil {
		state := vo.SessionState(*req.State)
		filter.State = &state
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	sessions, total, err := s.SessionRepo.List(ctx, filter, req.Offset, limit)
	if err != nil {
		return nil, err
	}
	return &dto.SessionListResponse{
		Sessions: assembler.ToSessionResponseList(sessions),
		Total:    total,
		Offset:   req.Offset,
		Limit:    limit,
	}, nil
}

func (s *agentAppService) DeleteSession(ctx context.Context, id vo.SessionID) error {
	return s.SessionRepo.Delete(ctx, id)
}
