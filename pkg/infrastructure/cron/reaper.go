package cron

import (
	"context"
	"regexp"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/conversation/repository"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

const (
	// minReaperIntervalMs is the minimum allowed reaper scan interval (5 minutes).
	minReaperIntervalMs = 300000
	// defaultRetentionMs is the default session retention period (24 hours).
	defaultRetentionMs = 86400000
	// reaperPageSize is the batch size for listing sessions during cleanup.
	reaperPageSize = 100
)

// cronRunSessionPattern matches temporary cron run session keys: cron:<jobId>:run:<uuid>
var cronRunSessionPattern = regexp.MustCompile(`^cron:[^:]+:run:[0-9a-f-]+$`)

// Reaper periodically cleans up expired isolated cron run sessions.
type Reaper struct {
	sessionRepo repository.SessionRepository
	intervalMs  int64
	retentionMs int64
	enabled     bool
	logger      logger.Logger
	cancel      context.CancelFunc
}

// NewReaper creates a new Reaper. If intervalMs is below the minimum (300000ms),
// it is clamped to the minimum. If retentionMs is zero, the default (24h) is used.
func NewReaper(
	sessionRepo repository.SessionRepository,
	intervalMs int64,
	retentionMs int64,
	enabled bool,
	log logger.Logger,
) *Reaper {
	if intervalMs < minReaperIntervalMs {
		intervalMs = minReaperIntervalMs
	}
	if retentionMs <= 0 {
		retentionMs = defaultRetentionMs
	}
	return &Reaper{
		sessionRepo: sessionRepo,
		intervalMs:  intervalMs,
		retentionMs: retentionMs,
		enabled:     enabled,
		logger:      log,
	}
}

// Start begins the periodic session cleanup loop. It blocks until the context
// is cancelled or Stop is called.
func (r *Reaper) Start(ctx context.Context) {
	if !r.enabled {
		r.logger.Info("cron session reaper disabled")
		return
	}

	ctx, r.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(time.Duration(r.intervalMs) * time.Millisecond)
		defer ticker.Stop()

		r.logger.Info("cron session reaper started",
			"intervalMs", r.intervalMs,
			"retentionMs", r.retentionMs,
		)

		for {
			select {
			case <-ctx.Done():
				r.logger.Info("cron session reaper stopping")
				return
			case <-ticker.C:
				r.sweep(ctx)
			}
		}
	}()
}

// Stop cancels the reaper's context, causing the sweep loop to exit.
func (r *Reaper) Stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

// sweep scans all sessions in batches, identifies expired cron run sessions,
// and deletes them.
func (r *Reaper) sweep(ctx context.Context) {
	retentionThreshold := time.Now().Add(-time.Duration(r.retentionMs) * time.Millisecond)
	pruned := 0
	offset := 0

	for {
		sessions, total, err := r.sessionRepo.List(ctx, repository.SessionFilter{}, offset, reaperPageSize)
		if err != nil {
			r.logger.Error("reaper: list sessions failed", "error", err)
			return
		}

		for _, session := range sessions {
			key := session.ID().Value()

			// Only target temporary cron run sessions (cron:<jobId>:run:<uuid>).
			// Base sessions (cron:<jobId>) are never cleaned up.
			if !cronRunSessionPattern.MatchString(key) {
				continue
			}

			// Check if the session has expired beyond the retention period.
			if session.LastActiveAt().Before(retentionThreshold) {
				if err := r.sessionRepo.Delete(ctx, session.ID()); err != nil {
					r.logger.Error("reaper: delete session failed",
						"sessionId", key,
						"error", err,
					)
					continue
				}
				pruned++
			}
		}

		offset += reaperPageSize
		if offset >= total {
			break
		}
	}

	if pruned > 0 {
		r.logger.Info("pruned expired cron run session(s)", "count", pruned)
	}
}
