package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/robfig/cron/v3"

	cron_domain "github.com/make-bin/groundhog/pkg/domain/cron"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
)

// cronParser supports both 5-field and 6-field (with optional seconds) expressions and descriptors.
var cronParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

type cronSchedulerServiceImpl struct{}

// ComputeNextRun computes the next run time in milliseconds for the given schedule.
// Returns nil if there is no next run (e.g., an "at" schedule that has already passed).
func (s *cronSchedulerServiceImpl) ComputeNextRun(schedule vo.CronSchedule, createdAtMs, nowMs int64) (*int64, error) {
	switch schedule.Kind() {
	case vo.ScheduleKindAt:
		return s.computeNextRunAt(schedule, nowMs)
	case vo.ScheduleKindEvery:
		return s.computeNextRunEvery(schedule, createdAtMs, nowMs)
	case vo.ScheduleKindCron:
		return s.computeNextRunCron(schedule, nowMs)
	default:
		return nil, fmt.Errorf("%w: unknown schedule kind %q", cron_domain.ErrInvalidSchedule, schedule.Kind())
	}
}

// IsValidSchedule validates that the schedule is well-formed.
func (s *cronSchedulerServiceImpl) IsValidSchedule(schedule vo.CronSchedule) error {
	switch schedule.Kind() {
	case vo.ScheduleKindAt:
		if _, err := time.Parse(time.RFC3339, schedule.At()); err != nil {
			return fmt.Errorf("%w: at field must be a valid RFC3339 timestamp: %v", cron_domain.ErrInvalidSchedule, err)
		}
		return nil
	case vo.ScheduleKindEvery:
		if schedule.EveryMs() <= 0 {
			return fmt.Errorf("%w: everyMs must be greater than 0", cron_domain.ErrInvalidSchedule)
		}
		return nil
	case vo.ScheduleKindCron:
		if _, err := cronParser.Parse(schedule.Expr()); err != nil {
			return fmt.Errorf("%w: invalid cron expression: %v", cron_domain.ErrInvalidSchedule, err)
		}
		return nil
	default:
		return fmt.Errorf("%w: unknown schedule kind %q", cron_domain.ErrInvalidSchedule, schedule.Kind())
	}
}

func (s *cronSchedulerServiceImpl) computeNextRunAt(schedule vo.CronSchedule, nowMs int64) (*int64, error) {
	t, err := time.Parse(time.RFC3339, schedule.At())
	if err != nil {
		return nil, fmt.Errorf("%w: at field must be a valid RFC3339 timestamp: %v", cron_domain.ErrInvalidSchedule, err)
	}
	atMs := t.UnixMilli()
	if atMs <= nowMs {
		// Already past — no next run.
		return nil, nil
	}
	return &atMs, nil
}

func (s *cronSchedulerServiceImpl) computeNextRunEvery(schedule vo.CronSchedule, createdAtMs, nowMs int64) (*int64, error) {
	everyMs := schedule.EveryMs()
	if everyMs <= 0 {
		return nil, fmt.Errorf("%w: everyMs must be greater than 0", cron_domain.ErrInvalidSchedule)
	}

	anchor := createdAtMs
	if schedule.AnchorMs() != nil {
		anchor = *schedule.AnchorMs()
	}

	// Find smallest n such that anchor + n*everyMs > nowMs
	// n = (nowMs - anchor) / everyMs + 1  (integer division)
	n := (nowMs-anchor)/everyMs + 1
	next := anchor + n*everyMs
	return &next, nil
}

func (s *cronSchedulerServiceImpl) computeNextRunCron(schedule vo.CronSchedule, nowMs int64) (*int64, error) {
	parsed, err := cronParser.Parse(schedule.Expr())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid cron expression: %v", cron_domain.ErrInvalidSchedule, err)
	}

	tz := schedule.Tz()
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid timezone %q: %v", cron_domain.ErrInvalidSchedule, tz, err)
	}

	now := time.Now().In(loc)
	next := parsed.Next(now)
	if next.IsZero() {
		return nil, nil
	}

	nextMs := next.UnixMilli()

	if schedule.StaggerMs() > 0 {
		//nolint:gosec // non-cryptographic random jitter is intentional
		jitter := rand.Int63n(schedule.StaggerMs())
		nextMs += jitter
	}

	return &nextMs, nil
}
