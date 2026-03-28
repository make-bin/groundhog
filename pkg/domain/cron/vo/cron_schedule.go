package vo

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// ScheduleKind represents the type of cron schedule.
type ScheduleKind string

const (
	ScheduleKindAt    ScheduleKind = "at"
	ScheduleKindEvery ScheduleKind = "every"
	ScheduleKindCron  ScheduleKind = "cron"
)

// CronSchedule is an immutable value object representing a scheduling rule.
// It supports three kinds: at (one-shot), every (fixed interval), and cron (expression-based).
type CronSchedule struct {
	kind      ScheduleKind
	at        string // ISO8601/RFC3339, valid when kind=at
	everyMs   int64  // millisecond interval, valid when kind=every
	anchorMs  *int64 // optional anchor epoch-ms, valid when kind=every
	expr      string // cron expression, valid when kind=cron
	tz        string // timezone name, optional when kind=cron (default UTC)
	staggerMs int64  // random jitter window ms, optional when kind=cron
}

// NewCronScheduleAt creates a one-shot schedule that fires at the given time.
// at must be a non-empty RFC3339 / ISO8601 string.
func NewCronScheduleAt(at string) (CronSchedule, error) {
	if at == "" {
		return CronSchedule{}, fmt.Errorf("at must not be empty")
	}
	if _, err := time.Parse(time.RFC3339, at); err != nil {
		return CronSchedule{}, fmt.Errorf("at must be a valid RFC3339/ISO8601 timestamp: %w", err)
	}
	return CronSchedule{kind: ScheduleKindAt, at: at}, nil
}

// NewCronScheduleEvery creates a fixed-interval schedule.
// everyMs must be > 0. anchorMs is optional.
func NewCronScheduleEvery(everyMs int64, anchorMs *int64) (CronSchedule, error) {
	if everyMs <= 0 {
		return CronSchedule{}, fmt.Errorf("everyMs must be greater than 0, got %d", everyMs)
	}
	return CronSchedule{
		kind:     ScheduleKindEvery,
		everyMs:  everyMs,
		anchorMs: anchorMs,
	}, nil
}

// cronParser is a shared parser that supports both 5-field and 6-field (with seconds) expressions
// as well as descriptors like @hourly.
// Using OptionalSecond so that both "* * * * *" (5-field) and "* * * * * *" (6-field) are accepted.
var cronParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// NewCronScheduleCron creates an expression-based schedule.
// expr must be non-empty and parseable by robfig/cron/v3.
// tz, if non-empty, must be a valid IANA timezone name.
// staggerMs must be >= 0.
func NewCronScheduleCron(expr, tz string, staggerMs int64) (CronSchedule, error) {
	if expr == "" {
		return CronSchedule{}, fmt.Errorf("cron expression must not be empty")
	}
	if _, err := cronParser.Parse(expr); err != nil {
		return CronSchedule{}, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}
	if tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			return CronSchedule{}, fmt.Errorf("invalid timezone %q: %w", tz, err)
		}
	}
	if staggerMs < 0 {
		return CronSchedule{}, fmt.Errorf("staggerMs must be >= 0, got %d", staggerMs)
	}
	return CronSchedule{
		kind:      ScheduleKindCron,
		expr:      expr,
		tz:        tz,
		staggerMs: staggerMs,
	}, nil
}

// Kind returns the schedule kind.
func (s CronSchedule) Kind() ScheduleKind { return s.kind }

// At returns the ISO8601 timestamp string (valid when Kind == ScheduleKindAt).
func (s CronSchedule) At() string { return s.at }

// EveryMs returns the interval in milliseconds (valid when Kind == ScheduleKindEvery).
func (s CronSchedule) EveryMs() int64 { return s.everyMs }

// AnchorMs returns the optional anchor epoch-ms (valid when Kind == ScheduleKindEvery).
func (s CronSchedule) AnchorMs() *int64 { return s.anchorMs }

// Expr returns the cron expression string (valid when Kind == ScheduleKindCron).
func (s CronSchedule) Expr() string { return s.expr }

// Tz returns the timezone name (valid when Kind == ScheduleKindCron; empty means UTC).
func (s CronSchedule) Tz() string { return s.tz }

// StaggerMs returns the random jitter window in milliseconds (valid when Kind == ScheduleKindCron).
func (s CronSchedule) StaggerMs() int64 { return s.staggerMs }
