package mapper

import (
	"fmt"
	"time"

	cron_repository "github.com/make-bin/groundhog/pkg/domain/cron/repository"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/infrastructure/persistence/po"
)

// DomainToCronRunLogPO converts a CronRunLog domain entity to a CronRunLogPO.
func DomainToCronRunLogPO(log *cron_repository.CronRunLog) *po.CronRunLogPO {
	return &po.CronRunLogPO{
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
		CreatedAt:      time.Now(),
	}
}

// CronRunLogPOToDomain converts a CronRunLogPO to a CronRunLog domain entity.
func CronRunLogPOToDomain(p *po.CronRunLogPO) (*cron_repository.CronRunLog, error) {
	jobID, err := vo.NewCronJobID(p.JobID)
	if err != nil {
		return nil, fmt.Errorf("reconstruct job_id: %w", err)
	}

	return &cron_repository.CronRunLog{
		ID:             p.ID,
		JobID:          jobID,
		Ts:             p.Ts,
		Action:         p.Action,
		Status:         p.Status,
		Error:          p.Error,
		Summary:        p.Summary,
		SessionID:      p.SessionID,
		RunAtMs:        p.RunAtMs,
		DurationMs:     p.DurationMs,
		NextRunAtMs:    p.NextRunAtMs,
		Model:          p.Model,
		Provider:       p.Provider,
		InputTokens:    p.InputTokens,
		OutputTokens:   p.OutputTokens,
		TotalTokens:    p.TotalTokens,
		DeliveryStatus: p.DeliveryStatus,
		DeliveryError:  p.DeliveryError,
	}, nil
}
