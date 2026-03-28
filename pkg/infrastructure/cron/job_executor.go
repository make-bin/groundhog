package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/application/service"
	"github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	conversation_vo "github.com/make-bin/groundhog/pkg/domain/conversation/vo"
	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// ExecuteResult holds the outcome of a cron job execution.
type ExecuteResult struct {
	Status    string // "ok" or "error"
	Error     string
	Summary   string
	SessionID string
	Model     string
	Provider  string
	Usage     TokenUsage
}

// TokenUsage records LLM token consumption.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// JobExecutor executes cron job payloads (agentTurn / systemEvent).
type JobExecutor struct {
	agentAppSvc   service.AgentAppService
	sessionRepo   repository.SessionRepository
	channelAppSvc service.ChannelAppService
	logger        logger.Logger
}

// NewJobExecutor creates a new JobExecutor with the required dependencies.
func NewJobExecutor(
	agentAppSvc service.AgentAppService,
	sessionRepo repository.SessionRepository,
	channelAppSvc service.ChannelAppService,
	log logger.Logger,
) *JobExecutor {
	return &JobExecutor{
		agentAppSvc:   agentAppSvc,
		sessionRepo:   sessionRepo,
		channelAppSvc: channelAppSvc,
		logger:        log,
	}
}

// Execute runs the cron job payload and returns the result.
// It supports agentTurn and systemEvent payload kinds with timeout cancellation.
func (e *JobExecutor) Execute(ctx context.Context, job *cron_job.CronJob) (*ExecuteResult, error) {
	payload := job.Payload()

	// Apply timeout if configured for agentTurn.
	if payload.Kind() == vo.PayloadKindAgentTurn && payload.TimeoutSeconds() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(payload.TimeoutSeconds())*time.Second)
		defer cancel()
	}

	switch payload.Kind() {
	case vo.PayloadKindAgentTurn:
		return e.executeAgentTurn(ctx, job)
	case vo.PayloadKindSystemEvent:
		return e.executeSystemEvent(ctx, job)
	default:
		return nil, fmt.Errorf("unsupported payload kind: %s", payload.Kind())
	}
}

// executeAgentTurn resolves the target session and triggers an agent turn.
func (e *JobExecutor) executeAgentTurn(ctx context.Context, job *cron_job.CronJob) (*ExecuteResult, error) {
	sessionKey, err := e.resolveSessionKey(ctx, job)
	if err != nil {
		return &ExecuteResult{
			Status: "error",
			Error:  fmt.Sprintf("resolve session: %v", err),
		}, nil
	}

	sessionID, err := conversation_vo.NewSessionID(sessionKey)
	if err != nil {
		return &ExecuteResult{
			Status: "error",
			Error:  fmt.Sprintf("invalid session ID: %v", err),
		}, nil
	}

	// Ensure the session exists; create it if needed for isolated/current targets.
	if err := e.ensureSession(ctx, job, sessionID); err != nil {
		return &ExecuteResult{
			Status: "error",
			Error:  fmt.Sprintf("ensure session: %v", err),
		}, nil
	}

	payload := job.Payload()
	turnResp, err := e.agentAppSvc.ExecuteTurn(ctx, sessionID, payload.Message())
	if err != nil {
		return &ExecuteResult{
			Status:    "error",
			Error:     err.Error(),
			SessionID: sessionKey,
		}, nil
	}

	return &ExecuteResult{
		Status:    "ok",
		Summary:   truncate(turnResp.Response, 500),
		SessionID: sessionKey,
		Model:     turnResp.ModelUsed,
	}, nil
}

// executeSystemEvent injects a system message into the main session.
func (e *JobExecutor) executeSystemEvent(ctx context.Context, job *cron_job.CronJob) (*ExecuteResult, error) {
	agentID := job.AgentId()
	if agentID == "" {
		agentID = "default"
	}

	// For systemEvent, the session target is always "main".
	// We use the agent's main session key.
	sessionKey := fmt.Sprintf("main:%s", agentID)
	sessionID, err := conversation_vo.NewSessionID(sessionKey)
	if err != nil {
		return &ExecuteResult{
			Status: "error",
			Error:  fmt.Sprintf("invalid session ID: %v", err),
		}, nil
	}

	// Inject the system event text as a user message to the main session.
	payload := job.Payload()
	_, execErr := e.agentAppSvc.ExecuteTurn(ctx, sessionID, payload.Text())
	if execErr != nil {
		return &ExecuteResult{
			Status:    "error",
			Error:     execErr.Error(),
			SessionID: sessionKey,
		}, nil
	}

	return &ExecuteResult{
		Status:    "ok",
		Summary:   "system event injected",
		SessionID: sessionKey,
	}, nil
}

// resolveSessionKey determines the session key based on the job's sessionTarget.
func (e *JobExecutor) resolveSessionKey(ctx context.Context, job *cron_job.CronJob) (string, error) {
	target := job.SessionTarget()

	switch {
	case target == "isolated":
		// Create a unique session key for each run: cron:<jobId>:run:<uuid>
		runID := vo.GenerateCronJobID() // reuse UUID generator
		return fmt.Sprintf("cron:%s:run:%s", job.ID().Value(), runID.Value()), nil

	case target == "current":
		// Reuse the last session key if available; otherwise create a new one.
		if job.SessionKey() != "" {
			return job.SessionKey(), nil
		}
		return fmt.Sprintf("cron:%s", job.ID().Value()), nil

	case strings.HasPrefix(target, "session:"):
		// Use the specified session key.
		key := strings.TrimPrefix(target, "session:")
		if key == "" {
			return "", fmt.Errorf("session target key must not be empty")
		}
		return key, nil

	case target == "main":
		agentID := job.AgentId()
		if agentID == "" {
			agentID = "default"
		}
		return fmt.Sprintf("main:%s", agentID), nil

	default:
		return "", fmt.Errorf("unknown session target: %s", target)
	}
}

// ensureSession creates a session if it doesn't exist yet (for isolated/current targets).
func (e *JobExecutor) ensureSession(ctx context.Context, job *cron_job.CronJob, sessionID conversation_vo.SessionID) error {
	_, err := e.sessionRepo.FindByID(ctx, sessionID)
	if err == nil {
		return nil // session already exists
	}

	// Session not found — create a new one.
	agentID := job.AgentId()
	if agentID == "" {
		agentID = "default"
	}

	payload := job.Payload()
	model := payload.Model()
	if model == "" {
		model = "default"
	}

	req := &dto.CreateSessionRequest{
		AgentID:      agentID,
		UserID:       "cron-scheduler",
		Provider:     "gemini",
		ModelName:    model,
		SystemPrompt: fmt.Sprintf("Cron job: %s", job.Name()),
	}

	_, createErr := e.agentAppSvc.CreateSession(ctx, req)
	if createErr != nil {
		return fmt.Errorf("create session for cron job %s: %w", job.ID().Value(), createErr)
	}
	return nil
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
