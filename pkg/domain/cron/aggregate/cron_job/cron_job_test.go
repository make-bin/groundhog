package cron_job_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/make-bin/groundhog/pkg/domain/cron"
	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
)

// helpers

func makeSchedule(t *testing.T) vo.CronSchedule {
	t.Helper()
	s, err := vo.NewCronScheduleCron("* * * * *", "", 0)
	require.NoError(t, err)
	return s
}

func makeSystemEventPayload(t *testing.T) vo.CronPayload {
	t.Helper()
	p, err := vo.NewCronPayloadSystemEvent("hello")
	require.NoError(t, err)
	return p
}

func makeAgentTurnPayload(t *testing.T) vo.CronPayload {
	t.Helper()
	p, err := vo.NewCronPayloadAgentTurn("do something", "", false, 0, false)
	require.NoError(t, err)
	return p
}

func baseParams(t *testing.T) cron_job.CreateCronJobParams {
	t.Helper()
	return cron_job.CreateCronJobParams{
		ID:            vo.GenerateCronJobID(),
		Name:          "test-job",
		Enabled:       true,
		CreatedAtMs:   time.Now().UnixMilli(),
		UpdatedAtMs:   time.Now().UnixMilli(),
		Schedule:      makeSchedule(t),
		SessionTarget: "main",
		WakeMode:      "next-heartbeat",
		Payload:       makeSystemEventPayload(t),
		State:         vo.NewCronJobState(),
	}
}

// --- Constructor validation ---

func TestNewCronJob_EmptyName(t *testing.T) {
	p := baseParams(t)
	p.Name = ""
	_, err := cron_job.NewCronJob(p)
	assert.EqualError(t, err, "name must not be empty")
}

func TestNewCronJob_InvalidSessionTarget(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "unknown"
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrInvalidSessionTarget)
}

func TestNewCronJob_ValidSessionTargets(t *testing.T) {
	targets := []string{"main", "isolated", "current", "session:abc", "session:my-key-123"}
	for _, target := range targets {
		p := baseParams(t)
		p.SessionTarget = target
		if target == "main" {
			p.Payload = makeSystemEventPayload(t)
		} else {
			p.Payload = makeAgentTurnPayload(t)
		}
		_, err := cron_job.NewCronJob(p)
		assert.NoError(t, err, "target=%s", target)
	}
}

func TestNewCronJob_SessionKeyEmptyAfterColon(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "session:"
	p.Payload = makeAgentTurnPayload(t)
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrInvalidSessionTarget)
}

func TestNewCronJob_MainRequiresSystemEvent(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "main"
	p.Payload = makeAgentTurnPayload(t)
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrPayloadSessionMismatch)
}

func TestNewCronJob_IsolatedRequiresAgentTurn(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "isolated"
	p.Payload = makeSystemEventPayload(t)
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrPayloadSessionMismatch)
}

func TestNewCronJob_CurrentRequiresAgentTurn(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "current"
	p.Payload = makeSystemEventPayload(t)
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrPayloadSessionMismatch)
}

func TestNewCronJob_SessionKeyRequiresAgentTurn(t *testing.T) {
	p := baseParams(t)
	p.SessionTarget = "session:mykey"
	p.Payload = makeSystemEventPayload(t)
	_, err := cron_job.NewCronJob(p)
	assert.ErrorIs(t, err, cron.ErrPayloadSessionMismatch)
}

func TestNewCronJob_WakeModeDefaultsToNextHeartbeat(t *testing.T) {
	p := baseParams(t)
	p.WakeMode = ""
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.Equal(t, "next-heartbeat", j.WakeMode())
}

func TestNewCronJob_WakeModePreservedWhenSet(t *testing.T) {
	p := baseParams(t)
	p.WakeMode = "now"
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.Equal(t, "now", j.WakeMode())
}

func TestNewCronJob_Success(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.Equal(t, p.Name, j.Name())
	assert.Equal(t, p.SessionTarget, j.SessionTarget())
	assert.True(t, j.Enabled())
}

// --- ReconstructCronJob ---

func TestReconstructCronJob_BypassesValidation(t *testing.T) {
	p := baseParams(t)
	p.Name = "" // would fail NewCronJob
	j := cron_job.ReconstructCronJob(p)
	assert.NotNil(t, j)
	assert.Equal(t, "", j.Name())
}

// --- Enable / Disable ---

func TestEnable(t *testing.T) {
	p := baseParams(t)
	p.Enabled = false
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	before := j.UpdatedAtMs()
	time.Sleep(time.Millisecond)
	j.Enable()
	assert.True(t, j.Enabled())
	assert.Greater(t, j.UpdatedAtMs(), before)
}

func TestDisable(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	before := j.UpdatedAtMs()
	time.Sleep(time.Millisecond)
	j.Disable()
	assert.False(t, j.Enabled())
	assert.Greater(t, j.UpdatedAtMs(), before)
}

// --- MarkRunning ---

func TestMarkRunning(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	nowMs := time.Now().UnixMilli()
	j.MarkRunning(nowMs)
	require.NotNil(t, j.State().RunningAtMs())
	assert.Equal(t, nowMs, *j.State().RunningAtMs())
}

// --- MarkCompleted ---

func TestMarkCompleted_OkResetsConsecutiveErrors(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	// Simulate prior errors
	j.UpdateState(j.State().WithConsecutiveErrors(3))
	j.MarkCompleted("ok", 100, nil, "")
	assert.Equal(t, 0, j.State().ConsecutiveErrors())
	assert.Equal(t, "ok", j.State().LastRunStatus())
	assert.Nil(t, j.State().RunningAtMs())
}

func TestMarkCompleted_ErrorIncrementsConsecutiveErrors(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	j.MarkCompleted("error", 50, nil, "something went wrong")
	assert.Equal(t, 1, j.State().ConsecutiveErrors())
	assert.Equal(t, "error", j.State().LastRunStatus())
	assert.Equal(t, "something went wrong", j.State().LastError())
	assert.Equal(t, int64(50), j.State().LastDurationMs())
}

func TestMarkCompleted_SetsNextRunAtMs(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	next := time.Now().Add(time.Minute).UnixMilli()
	j.MarkCompleted("ok", 10, &next, "")
	require.NotNil(t, j.State().NextRunAtMs())
	assert.Equal(t, next, *j.State().NextRunAtMs())
}

// --- UpdateState ---

func TestUpdateState(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	newState := vo.NewCronJobState().WithConsecutiveErrors(5)
	j.UpdateState(newState)
	assert.Equal(t, 5, j.State().ConsecutiveErrors())
}

// --- ApplyPatch ---

func TestApplyPatch_UpdatesName(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	newName := "updated-name"
	err = j.ApplyPatch(cron_job.UpdateCronJobParams{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, "updated-name", j.Name())
}

func TestApplyPatch_InvalidSessionTarget(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	bad := "bad-target"
	err = j.ApplyPatch(cron_job.UpdateCronJobParams{SessionTarget: &bad})
	assert.ErrorIs(t, err, cron.ErrInvalidSessionTarget)
}

func TestApplyPatch_PayloadMismatch(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	// job is "main" + systemEvent; try to switch payload to agentTurn
	agentPayload := makeAgentTurnPayload(t)
	err = j.ApplyPatch(cron_job.UpdateCronJobParams{Payload: &agentPayload})
	assert.ErrorIs(t, err, cron.ErrPayloadSessionMismatch)
}

func TestApplyPatch_ChangeSessionTargetAndPayloadTogether(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	newTarget := "isolated"
	agentPayload := makeAgentTurnPayload(t)
	err = j.ApplyPatch(cron_job.UpdateCronJobParams{
		SessionTarget: &newTarget,
		Payload:       &agentPayload,
	})
	require.NoError(t, err)
	assert.Equal(t, "isolated", j.SessionTarget())
	assert.Equal(t, vo.PayloadKindAgentTurn, j.Payload().Kind())
}

func TestApplyPatch_UpdatesUpdatedAtMs(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	before := j.UpdatedAtMs()
	time.Sleep(time.Millisecond)
	newName := "x"
	err = j.ApplyPatch(cron_job.UpdateCronJobParams{Name: &newName})
	require.NoError(t, err)
	assert.Greater(t, j.UpdatedAtMs(), before)
}

// --- ShouldAlert ---

func TestShouldAlert_NilFailureAlert(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.False(t, j.ShouldAlert())
}

func TestShouldAlert_BelowThreshold(t *testing.T) {
	p := baseParams(t)
	alert, err := vo.NewCronFailureAlert(3, "", "", 0, "announce", "")
	require.NoError(t, err)
	p.FailureAlert = &alert
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	j.UpdateState(j.State().WithConsecutiveErrors(2))
	assert.False(t, j.ShouldAlert())
}

func TestShouldAlert_AtThresholdNoLastAlert(t *testing.T) {
	p := baseParams(t)
	alert, err := vo.NewCronFailureAlert(3, "", "", 0, "announce", "")
	require.NoError(t, err)
	p.FailureAlert = &alert
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	j.UpdateState(j.State().WithConsecutiveErrors(3))
	assert.True(t, j.ShouldAlert())
}

func TestShouldAlert_CooldownNotElapsed(t *testing.T) {
	p := baseParams(t)
	cooldown := int64(60_000) // 1 minute
	alert, err := vo.NewCronFailureAlert(3, "", "", cooldown, "announce", "")
	require.NoError(t, err)
	p.FailureAlert = &alert
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	recentMs := time.Now().UnixMilli() - 1000 // 1 second ago
	j.UpdateState(j.State().WithConsecutiveErrors(3).WithLastFailureAlertAtMs(&recentMs))
	assert.False(t, j.ShouldAlert())
}

func TestShouldAlert_CooldownElapsed(t *testing.T) {
	p := baseParams(t)
	cooldown := int64(60_000) // 1 minute
	alert, err := vo.NewCronFailureAlert(3, "", "", cooldown, "announce", "")
	require.NoError(t, err)
	p.FailureAlert = &alert
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	oldMs := time.Now().UnixMilli() - 120_000 // 2 minutes ago
	j.UpdateState(j.State().WithConsecutiveErrors(3).WithLastFailureAlertAtMs(&oldMs))
	assert.True(t, j.ShouldAlert())
}

// --- ShouldDeleteAfterRun ---

func TestShouldDeleteAfterRun_Nil(t *testing.T) {
	p := baseParams(t)
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.False(t, j.ShouldDeleteAfterRun())
}

func TestShouldDeleteAfterRun_False(t *testing.T) {
	p := baseParams(t)
	f := false
	p.DeleteAfterRun = &f
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.False(t, j.ShouldDeleteAfterRun())
}

func TestShouldDeleteAfterRun_True(t *testing.T) {
	p := baseParams(t)
	tr := true
	p.DeleteAfterRun = &tr
	j, err := cron_job.NewCronJob(p)
	require.NoError(t, err)
	assert.True(t, j.ShouldDeleteAfterRun())
}
