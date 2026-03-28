package assembler

import (
	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
)

// ToCronJobResponse converts a CronJob domain aggregate to a CronJobResponse DTO.
func ToCronJobResponse(job *cron_job.CronJob) *dto.CronJobResponse {
	if job == nil {
		return nil
	}

	sched := job.Schedule()
	schedDTO := dto.CronScheduleDTO{
		Kind:      string(sched.Kind()),
		At:        sched.At(),
		EveryMs:   sched.EveryMs(),
		AnchorMs:  sched.AnchorMs(),
		Expr:      sched.Expr(),
		Tz:        sched.Tz(),
		StaggerMs: sched.StaggerMs(),
	}

	payload := job.Payload()
	payloadDTO := dto.CronPayloadDTO{
		Kind:           string(payload.Kind()),
		Text:           payload.Text(),
		Message:        payload.Message(),
		Model:          payload.Model(),
		Thinking:       payload.Thinking(),
		TimeoutSeconds: payload.TimeoutSeconds(),
		LightContext:   payload.LightContext(),
	}

	var deliveryDTO *dto.CronDeliveryDTO
	if d := job.Delivery(); d != nil {
		deliveryDTO = &dto.CronDeliveryDTO{
			Mode:               string(d.Mode()),
			Channel:            d.Channel(),
			To:                 d.To(),
			AccountId:          d.AccountId(),
			BestEffort:         d.BestEffort(),
			FailureDestination: d.FailureDestination(),
		}
	}

	var failureAlertDTO *dto.CronFailureAlertDTO
	if fa := job.FailureAlert(); fa != nil {
		failureAlertDTO = &dto.CronFailureAlertDTO{
			After:      fa.After(),
			Channel:    fa.Channel(),
			To:         fa.To(),
			CooldownMs: fa.CooldownMs(),
			Mode:       fa.Mode(),
			AccountId:  fa.AccountId(),
		}
	}

	state := job.State()
	stateDTO := dto.CronJobStateDTO{
		NextRunAtMs:          state.NextRunAtMs(),
		RunningAtMs:          state.RunningAtMs(),
		LastRunAtMs:          state.LastRunAtMs(),
		LastRunStatus:        state.LastRunStatus(),
		LastError:            state.LastError(),
		LastDurationMs:       state.LastDurationMs(),
		ConsecutiveErrors:    state.ConsecutiveErrors(),
		LastFailureAlertAtMs: state.LastFailureAlertAtMs(),
		ScheduleErrorCount:   state.ScheduleErrorCount(),
		LastDeliveryStatus:   state.LastDeliveryStatus(),
		LastDeliveryError:    state.LastDeliveryError(),
	}

	return &dto.CronJobResponse{
		ID:             job.ID().Value(),
		AgentId:        job.AgentId(),
		SessionKey:     job.SessionKey(),
		Name:           job.Name(),
		Description:    job.Description(),
		Enabled:        job.Enabled(),
		DeleteAfterRun: job.DeleteAfterRun(),
		CreatedAtMs:    job.CreatedAtMs(),
		UpdatedAtMs:    job.UpdatedAtMs(),
		Schedule:       schedDTO,
		SessionTarget:  job.SessionTarget(),
		WakeMode:       job.WakeMode(),
		Payload:        payloadDTO,
		Delivery:       deliveryDTO,
		FailureAlert:   failureAlertDTO,
		State:          stateDTO,
	}
}

// ToCronJobListResponse converts a slice of CronJob aggregates to a CronJobListResponse DTO.
func ToCronJobListResponse(jobs []*cron_job.CronJob, total, offset, limit int) *dto.CronJobListResponse {
	responses := make([]*dto.CronJobResponse, 0, len(jobs))
	for _, j := range jobs {
		responses = append(responses, ToCronJobResponse(j))
	}
	return &dto.CronJobListResponse{
		Jobs:   responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}
}

// ToRunLogResponse converts a CronRunLog entity to a CronRunLogResponse DTO.
func ToRunLogResponse(log *cron_repository.CronRunLog) *dto.CronRunLogResponse {
	if log == nil {
		return nil
	}
	return &dto.CronRunLogResponse{
		ID:             log.ID,
		JobID:          log.JobID.Value(),
		Ts:             log.Ts,
		Action:         log.Action,
		Status:         log.Status,
		Error:          log.Error,
		Summary:        log.Summary,
		SessionID:      log.SessionID,
		RunAtMs:        log.RunAtMs,
		DurationMs:     log.DurationMs,
		NextRunAtMs:    log.NextRunAtMs,
		Model:          log.Model,
		Provider:       log.Provider,
		InputTokens:    log.InputTokens,
		OutputTokens:   log.OutputTokens,
		TotalTokens:    log.TotalTokens,
		DeliveryStatus: log.DeliveryStatus,
		DeliveryError:  log.DeliveryError,
	}
}

// ToRunLogListResponse converts a slice of CronRunLog entities to a RunLogListResponse DTO.
func ToRunLogListResponse(logs []*cron_repository.CronRunLog, total, offset, limit int) *dto.RunLogListResponse {
	responses := make([]*dto.CronRunLogResponse, 0, len(logs))
	for _, l := range logs {
		responses = append(responses, ToRunLogResponse(l))
	}
	return &dto.RunLogListResponse{
		Logs:   responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}
}
