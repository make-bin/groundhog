package mapper

import (
	"encoding/json"
	"fmt"

	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// --- JSON intermediate structs ---

type scheduleJSON struct {
	Kind      string `json:"kind"`
	At        string `json:"at,omitempty"`
	EveryMs   int64  `json:"every_ms,omitempty"`
	AnchorMs  *int64 `json:"anchor_ms,omitempty"`
	Expr      string `json:"expr,omitempty"`
	Tz        string `json:"tz,omitempty"`
	StaggerMs int64  `json:"stagger_ms,omitempty"`
}

type payloadJSON struct {
	Kind           string `json:"kind"`
	Text           string `json:"text,omitempty"`
	Message        string `json:"message,omitempty"`
	Model          string `json:"model,omitempty"`
	Thinking       bool   `json:"thinking,omitempty"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
	LightContext   bool   `json:"light_context,omitempty"`
}

type deliveryJSON struct {
	Mode               string `json:"mode"`
	Channel            string `json:"channel,omitempty"`
	To                 string `json:"to,omitempty"`
	AccountId          string `json:"account_id,omitempty"`
	BestEffort         bool   `json:"best_effort,omitempty"`
	FailureDestination string `json:"failure_destination,omitempty"`
}

type failureAlertJSON struct {
	After      int    `json:"after"`
	Channel    string `json:"channel,omitempty"`
	To         string `json:"to,omitempty"`
	CooldownMs int64  `json:"cooldown_ms,omitempty"`
	Mode       string `json:"mode"`
	AccountId  string `json:"account_id,omitempty"`
}

type stateJSON struct {
	NextRunAtMs          *int64 `json:"next_run_at_ms"`
	RunningAtMs          *int64 `json:"running_at_ms"`
	LastRunAtMs          *int64 `json:"last_run_at_ms"`
	LastRunStatus        string `json:"last_run_status"`
	LastError            string `json:"last_error"`
	LastDurationMs       int64  `json:"last_duration_ms"`
	ConsecutiveErrors    int    `json:"consecutive_errors"`
	LastFailureAlertAtMs *int64 `json:"last_failure_alert_at_ms"`
	ScheduleErrorCount   int    `json:"schedule_error_count"`
	LastDeliveryStatus   string `json:"last_delivery_status"`
	LastDeliveryError    string `json:"last_delivery_error"`
}

// DomainToCronJobPO converts a CronJob aggregate to a CronJobPO.
func DomainToCronJobPO(job *cron_job.CronJob) (*po.CronJobPO, error) {
	// Serialize schedule
	sched := job.Schedule()
	schedBytes, err := json.Marshal(scheduleJSON{
		Kind:      string(sched.Kind()),
		At:        sched.At(),
		EveryMs:   sched.EveryMs(),
		AnchorMs:  sched.AnchorMs(),
		Expr:      sched.Expr(),
		Tz:        sched.Tz(),
		StaggerMs: sched.StaggerMs(),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal schedule: %w", err)
	}

	// Serialize payload
	pl := job.Payload()
	payloadBytes, err := json.Marshal(payloadJSON{
		Kind:           string(pl.Kind()),
		Text:           pl.Text(),
		Message:        pl.Message(),
		Model:          pl.Model(),
		Thinking:       pl.Thinking(),
		TimeoutSeconds: pl.TimeoutSeconds(),
		LightContext:   pl.LightContext(),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	// Serialize delivery (optional)
	deliveryStr := ""
	if d := job.Delivery(); d != nil {
		deliveryBytes, err := json.Marshal(deliveryJSON{
			Mode:               string(d.Mode()),
			Channel:            d.Channel(),
			To:                 d.To(),
			AccountId:          d.AccountId(),
			BestEffort:         d.BestEffort(),
			FailureDestination: d.FailureDestination(),
		})
		if err != nil {
			return nil, fmt.Errorf("marshal delivery: %w", err)
		}
		deliveryStr = string(deliveryBytes)
	}

	// Serialize failureAlert (optional)
	failureAlertStr := ""
	if fa := job.FailureAlert(); fa != nil {
		faBytes, err := json.Marshal(failureAlertJSON{
			After:      fa.After(),
			Channel:    fa.Channel(),
			To:         fa.To(),
			CooldownMs: fa.CooldownMs(),
			Mode:       fa.Mode(),
			AccountId:  fa.AccountId(),
		})
		if err != nil {
			return nil, fmt.Errorf("marshal failure_alert: %w", err)
		}
		failureAlertStr = string(faBytes)
	}

	// Serialize state
	st := job.State()
	stateBytes, err := json.Marshal(stateJSON{
		NextRunAtMs:          st.NextRunAtMs(),
		RunningAtMs:          st.RunningAtMs(),
		LastRunAtMs:          st.LastRunAtMs(),
		LastRunStatus:        st.LastRunStatus(),
		LastError:            st.LastError(),
		LastDurationMs:       st.LastDurationMs(),
		ConsecutiveErrors:    st.ConsecutiveErrors(),
		LastFailureAlertAtMs: st.LastFailureAlertAtMs(),
		ScheduleErrorCount:   st.ScheduleErrorCount(),
		LastDeliveryStatus:   st.LastDeliveryStatus(),
		LastDeliveryError:    st.LastDeliveryError(),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal state: %w", err)
	}

	// Use "null" for optional JSONB fields when empty, so PostgreSQL accepts them.
	if deliveryStr == "" {
		deliveryStr = "null"
	}
	if failureAlertStr == "" {
		failureAlertStr = "null"
	}

	return &po.CronJobPO{
		JobID:          job.ID().Value(),
		AgentId:        job.AgentId(),
		SessionKey:     job.SessionKey(),
		Name:           job.Name(),
		Description:    job.Description(),
		Enabled:        job.Enabled(),
		DeleteAfterRun: job.DeleteAfterRun(),
		Schedule:       string(schedBytes),
		SessionTarget:  job.SessionTarget(),
		WakeMode:       job.WakeMode(),
		Payload:        string(payloadBytes),
		Delivery:       deliveryStr,
		FailureAlert:   failureAlertStr,
		State:          string(stateBytes),
	}, nil
}

// CronJobPOToDomain converts a CronJobPO to a CronJob aggregate.
func CronJobPOToDomain(p *po.CronJobPO) (*cron_job.CronJob, error) {
	// Reconstruct CronJobID
	jobID, err := vo.NewCronJobID(p.JobID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct job_id: %w", err)
	}

	// Deserialize schedule
	var sj scheduleJSON
	if err := json.Unmarshal([]byte(p.Schedule), &sj); err != nil {
		return nil, fmt.Errorf("unmarshal schedule: %w", err)
	}
	var schedule vo.CronSchedule
	switch vo.ScheduleKind(sj.Kind) {
	case vo.ScheduleKindAt:
		schedule, err = vo.NewCronScheduleAt(sj.At)
	case vo.ScheduleKindEvery:
		schedule, err = vo.NewCronScheduleEvery(sj.EveryMs, sj.AnchorMs)
	case vo.ScheduleKindCron:
		schedule, err = vo.NewCronScheduleCron(sj.Expr, sj.Tz, sj.StaggerMs)
	default:
		return nil, fmt.Errorf("unknown schedule kind: %q", sj.Kind)
	}
	if err != nil {
		return nil, fmt.Errorf("reconstruct schedule: %w", err)
	}

	// Deserialize payload
	var pj payloadJSON
	if err := json.Unmarshal([]byte(p.Payload), &pj); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}
	var payload vo.CronPayload
	switch vo.PayloadKind(pj.Kind) {
	case vo.PayloadKindSystemEvent:
		payload, err = vo.NewCronPayloadSystemEvent(pj.Text)
	case vo.PayloadKindAgentTurn:
		payload, err = vo.NewCronPayloadAgentTurn(pj.Message, pj.Model, pj.Thinking, pj.TimeoutSeconds, pj.LightContext)
	default:
		return nil, fmt.Errorf("unknown payload kind: %q", pj.Kind)
	}
	if err != nil {
		return nil, fmt.Errorf("reconstruct payload: %w", err)
	}

	// Deserialize delivery (optional)
	var delivery *vo.CronDelivery
	if p.Delivery != "" && p.Delivery != "{}" && p.Delivery != "null" {
		var dj deliveryJSON
		if err := json.Unmarshal([]byte(p.Delivery), &dj); err != nil {
			return nil, fmt.Errorf("unmarshal delivery: %w", err)
		}
		d, err := vo.NewCronDelivery(vo.DeliveryMode(dj.Mode), dj.Channel, dj.To, dj.AccountId, dj.BestEffort, dj.FailureDestination)
		if err != nil {
			return nil, fmt.Errorf("reconstruct delivery: %w", err)
		}
		delivery = &d
	}

	// Deserialize failureAlert (optional)
	var failureAlert *vo.CronFailureAlert
	if p.FailureAlert != "" && p.FailureAlert != "{}" && p.FailureAlert != "null" {
		var faj failureAlertJSON
		if err := json.Unmarshal([]byte(p.FailureAlert), &faj); err != nil {
			return nil, fmt.Errorf("unmarshal failure_alert: %w", err)
		}
		fa, err := vo.NewCronFailureAlert(faj.After, faj.Channel, faj.To, faj.CooldownMs, faj.Mode, faj.AccountId)
		if err != nil {
			return nil, fmt.Errorf("reconstruct failure_alert: %w", err)
		}
		failureAlert = &fa
	}

	// Deserialize state
	var sj2 stateJSON
	if p.State != "" && p.State != "{}" {
		if err := json.Unmarshal([]byte(p.State), &sj2); err != nil {
			return nil, fmt.Errorf("unmarshal state: %w", err)
		}
	}
	state := vo.NewCronJobState().
		WithNextRunAtMs(sj2.NextRunAtMs).
		WithRunningAtMs(sj2.RunningAtMs).
		WithLastRunAtMs(sj2.LastRunAtMs).
		WithLastRunStatus(sj2.LastRunStatus).
		WithLastError(sj2.LastError).
		WithLastDurationMs(sj2.LastDurationMs).
		WithConsecutiveErrors(sj2.ConsecutiveErrors).
		WithLastFailureAlertAtMs(sj2.LastFailureAlertAtMs).
		WithScheduleErrorCount(sj2.ScheduleErrorCount).
		WithLastDeliveryStatus(sj2.LastDeliveryStatus).
		WithLastDeliveryError(sj2.LastDeliveryError)

	return cron_job.ReconstructCronJob(cron_job.CreateCronJobParams{
		ID:             jobID,
		AgentId:        p.AgentId,
		SessionKey:     p.SessionKey,
		Name:           p.Name,
		Description:    p.Description,
		Enabled:        p.Enabled,
		DeleteAfterRun: p.DeleteAfterRun,
		CreatedAtMs:    p.CreatedAt.UnixMilli(),
		UpdatedAtMs:    p.UpdatedAt.UnixMilli(),
		Schedule:       schedule,
		SessionTarget:  p.SessionTarget,
		WakeMode:       p.WakeMode,
		Payload:        payload,
		Delivery:       delivery,
		FailureAlert:   failureAlert,
		State:          state,
	}), nil
}
