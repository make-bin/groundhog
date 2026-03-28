package cron_job

import (
	"fmt"
	"regexp"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/cron"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
)

// sessionKeyPattern matches "session:<key>" where key is non-empty.
var sessionKeyPattern = regexp.MustCompile(`^session:.+$`)

// CronJob is the aggregate root for a scheduled cron task.
type CronJob struct {
	id             vo.CronJobID
	agentId        string
	sessionKey     string
	name           string
	description    string
	enabled        bool
	deleteAfterRun *bool
	createdAtMs    int64
	updatedAtMs    int64
	schedule       vo.CronSchedule
	sessionTarget  string // main/isolated/current/session:<key>
	wakeMode       string // now/next-heartbeat
	payload        vo.CronPayload
	delivery       *vo.CronDelivery
	failureAlert   *vo.CronFailureAlert
	state          vo.CronJobState
}

// CreateCronJobParams holds all parameters required to create a CronJob.
type CreateCronJobParams struct {
	ID             vo.CronJobID
	AgentId        string
	SessionKey     string
	Name           string
	Description    string
	Enabled        bool
	DeleteAfterRun *bool
	CreatedAtMs    int64
	UpdatedAtMs    int64
	Schedule       vo.CronSchedule
	SessionTarget  string
	WakeMode       string
	Payload        vo.CronPayload
	Delivery       *vo.CronDelivery
	FailureAlert   *vo.CronFailureAlert
	State          vo.CronJobState
}

// UpdateCronJobParams holds optional fields for patching a CronJob.
type UpdateCronJobParams struct {
	Name           *string
	Description    *string
	Enabled        *bool
	DeleteAfterRun *bool
	Schedule       *vo.CronSchedule
	SessionTarget  *string
	WakeMode       *string
	Payload        *vo.CronPayload
	Delivery       *vo.CronDelivery
	FailureAlert   *vo.CronFailureAlert
}

// isValidSessionTarget returns true if the given target is one of the allowed values.
func isValidSessionTarget(target string) bool {
	switch target {
	case "main", "isolated", "current":
		return true
	}
	return sessionKeyPattern.MatchString(target)
}

// validatePayloadSessionMatch checks that the payload kind is compatible with the session target.
func validatePayloadSessionMatch(sessionTarget string, payload vo.CronPayload) error {
	if sessionTarget == "main" {
		if payload.Kind() != vo.PayloadKindSystemEvent {
			return cron.ErrPayloadSessionMismatch
		}
		return nil
	}
	// isolated, current, or session:<key>
	if payload.Kind() != vo.PayloadKindAgentTurn {
		return cron.ErrPayloadSessionMismatch
	}
	return nil
}

// NewCronJob creates and validates a new CronJob aggregate root.
func NewCronJob(params CreateCronJobParams) (*CronJob, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	if !isValidSessionTarget(params.SessionTarget) {
		return nil, cron.ErrInvalidSessionTarget
	}

	if err := validatePayloadSessionMatch(params.SessionTarget, params.Payload); err != nil {
		return nil, err
	}

	wakeMode := params.WakeMode
	if wakeMode == "" {
		wakeMode = "next-heartbeat"
	}

	return &CronJob{
		id:             params.ID,
		agentId:        params.AgentId,
		sessionKey:     params.SessionKey,
		name:           params.Name,
		description:    params.Description,
		enabled:        params.Enabled,
		deleteAfterRun: params.DeleteAfterRun,
		createdAtMs:    params.CreatedAtMs,
		updatedAtMs:    params.UpdatedAtMs,
		schedule:       params.Schedule,
		sessionTarget:  params.SessionTarget,
		wakeMode:       wakeMode,
		payload:        params.Payload,
		delivery:       params.Delivery,
		failureAlert:   params.FailureAlert,
		state:          params.State,
	}, nil
}

// ReconstructCronJob rebuilds a CronJob from persisted data, bypassing validation.
// This should only be used by repository implementations.
func ReconstructCronJob(params CreateCronJobParams) *CronJob {
	wakeMode := params.WakeMode
	if wakeMode == "" {
		wakeMode = "next-heartbeat"
	}
	return &CronJob{
		id:             params.ID,
		agentId:        params.AgentId,
		sessionKey:     params.SessionKey,
		name:           params.Name,
		description:    params.Description,
		enabled:        params.Enabled,
		deleteAfterRun: params.DeleteAfterRun,
		createdAtMs:    params.CreatedAtMs,
		updatedAtMs:    params.UpdatedAtMs,
		schedule:       params.Schedule,
		sessionTarget:  params.SessionTarget,
		wakeMode:       wakeMode,
		payload:        params.Payload,
		delivery:       params.Delivery,
		failureAlert:   params.FailureAlert,
		state:          params.State,
	}
}

// --- Getters ---

func (j *CronJob) ID() vo.CronJobID                   { return j.id }
func (j *CronJob) AgentId() string                    { return j.agentId }
func (j *CronJob) SessionKey() string                 { return j.sessionKey }
func (j *CronJob) Name() string                       { return j.name }
func (j *CronJob) Description() string                { return j.description }
func (j *CronJob) Enabled() bool                      { return j.enabled }
func (j *CronJob) DeleteAfterRun() *bool              { return j.deleteAfterRun }
func (j *CronJob) CreatedAtMs() int64                 { return j.createdAtMs }
func (j *CronJob) UpdatedAtMs() int64                 { return j.updatedAtMs }
func (j *CronJob) Schedule() vo.CronSchedule          { return j.schedule }
func (j *CronJob) SessionTarget() string              { return j.sessionTarget }
func (j *CronJob) WakeMode() string                   { return j.wakeMode }
func (j *CronJob) Payload() vo.CronPayload            { return j.payload }
func (j *CronJob) Delivery() *vo.CronDelivery         { return j.delivery }
func (j *CronJob) FailureAlert() *vo.CronFailureAlert { return j.failureAlert }
func (j *CronJob) State() vo.CronJobState             { return j.state }

// --- Business Methods ---

// Enable sets the job as enabled and updates updatedAtMs.
func (j *CronJob) Enable() {
	j.enabled = true
	j.updatedAtMs = time.Now().UnixMilli()
}

// Disable sets the job as disabled and updates updatedAtMs.
func (j *CronJob) Disable() {
	j.enabled = false
	j.updatedAtMs = time.Now().UnixMilli()
}

// MarkRunning records that the job has started running at nowMs.
func (j *CronJob) MarkRunning(nowMs int64) {
	j.state = j.state.WithRunningAtMs(&nowMs)
}

// MarkCompleted clears the running state and records the result of a run.
// If status == "ok", consecutiveErrors is reset to 0; otherwise it is incremented.
func (j *CronJob) MarkCompleted(status string, durationMs int64, nextRunAtMs *int64, errMsg string) {
	s := j.state.WithRunningAtMs(nil)

	nowMs := time.Now().UnixMilli()
	s = s.WithLastRunAtMs(&nowMs)
	s = s.WithLastRunStatus(status)
	s = s.WithLastDurationMs(durationMs)
	s = s.WithLastError(errMsg)
	s = s.WithNextRunAtMs(nextRunAtMs)

	if status == "ok" {
		s = s.WithConsecutiveErrors(0)
	} else {
		s = s.WithConsecutiveErrors(s.ConsecutiveErrors() + 1)
	}

	j.state = s
}

// UpdateState replaces the entire state value object.
func (j *CronJob) UpdateState(state vo.CronJobState) {
	j.state = state
}

// ApplyPatch applies non-nil fields from patch to the job.
// Validates sessionTarget and payload constraints if either changes.
func (j *CronJob) ApplyPatch(patch UpdateCronJobParams) error {
	// Determine effective sessionTarget and payload after patch
	effectiveTarget := j.sessionTarget
	effectivePayload := j.payload

	if patch.SessionTarget != nil {
		effectiveTarget = *patch.SessionTarget
	}
	if patch.Payload != nil {
		effectivePayload = *patch.Payload
	}

	// Validate sessionTarget if it changed
	if patch.SessionTarget != nil {
		if !isValidSessionTarget(effectiveTarget) {
			return cron.ErrInvalidSessionTarget
		}
	}

	// Validate payload/session compatibility if either changed
	if patch.SessionTarget != nil || patch.Payload != nil {
		if err := validatePayloadSessionMatch(effectiveTarget, effectivePayload); err != nil {
			return err
		}
	}

	// Apply all non-nil fields
	if patch.Name != nil {
		j.name = *patch.Name
	}
	if patch.Description != nil {
		j.description = *patch.Description
	}
	if patch.Enabled != nil {
		j.enabled = *patch.Enabled
	}
	if patch.DeleteAfterRun != nil {
		j.deleteAfterRun = patch.DeleteAfterRun
	}
	if patch.Schedule != nil {
		j.schedule = *patch.Schedule
	}
	if patch.SessionTarget != nil {
		j.sessionTarget = effectiveTarget
	}
	if patch.WakeMode != nil {
		j.wakeMode = *patch.WakeMode
	}
	if patch.Payload != nil {
		j.payload = effectivePayload
	}
	if patch.Delivery != nil {
		j.delivery = patch.Delivery
	}
	if patch.FailureAlert != nil {
		j.failureAlert = patch.FailureAlert
	}

	j.updatedAtMs = time.Now().UnixMilli()
	return nil
}

// ShouldAlert returns true if a failure alert should be sent now.
// Conditions: failureAlert is configured, consecutiveErrors >= failureAlert.After(),
// and either no alert has been sent yet or the cooldown period has elapsed.
func (j *CronJob) ShouldAlert() bool {
	if j.failureAlert == nil {
		return false
	}
	if j.state.ConsecutiveErrors() < j.failureAlert.After() {
		return false
	}

	lastAlertAtMs := j.state.LastFailureAlertAtMs()
	if lastAlertAtMs == nil {
		return true
	}

	cooldownMs := j.failureAlert.CooldownMs()
	if cooldownMs <= 0 {
		return true
	}

	nowMs := time.Now().UnixMilli()
	return nowMs-*lastAlertAtMs >= cooldownMs
}

// ShouldDeleteAfterRun returns true if the job should be deleted after a successful run.
func (j *CronJob) ShouldDeleteAfterRun() bool {
	return j.deleteAfterRun != nil && *j.deleteAfterRun
}
