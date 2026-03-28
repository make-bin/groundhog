package vo

// CronJobState holds the mutable runtime state of a cron job.
// Mutations return a new copy (copy-on-write style).
type CronJobState struct {
	nextRunAtMs          *int64
	runningAtMs          *int64
	lastRunAtMs          *int64
	lastRunStatus        string
	lastError            string
	lastDurationMs       int64
	consecutiveErrors    int
	lastFailureAlertAtMs *int64
	scheduleErrorCount   int
	lastDeliveryStatus   string
	lastDeliveryError    string
}

// NewCronJobState returns a zero-value CronJobState.
func NewCronJobState() CronJobState {
	return CronJobState{}
}

// Getters

func (s CronJobState) NextRunAtMs() *int64          { return s.nextRunAtMs }
func (s CronJobState) RunningAtMs() *int64          { return s.runningAtMs }
func (s CronJobState) LastRunAtMs() *int64          { return s.lastRunAtMs }
func (s CronJobState) LastRunStatus() string        { return s.lastRunStatus }
func (s CronJobState) LastError() string            { return s.lastError }
func (s CronJobState) LastDurationMs() int64        { return s.lastDurationMs }
func (s CronJobState) ConsecutiveErrors() int       { return s.consecutiveErrors }
func (s CronJobState) LastFailureAlertAtMs() *int64 { return s.lastFailureAlertAtMs }
func (s CronJobState) ScheduleErrorCount() int      { return s.scheduleErrorCount }
func (s CronJobState) LastDeliveryStatus() string   { return s.lastDeliveryStatus }
func (s CronJobState) LastDeliveryError() string    { return s.lastDeliveryError }

// With* methods return a new CronJobState with the specified field updated.

func (s CronJobState) WithNextRunAtMs(v *int64) CronJobState {
	s.nextRunAtMs = v
	return s
}

func (s CronJobState) WithRunningAtMs(v *int64) CronJobState {
	s.runningAtMs = v
	return s
}

func (s CronJobState) WithLastRunAtMs(v *int64) CronJobState {
	s.lastRunAtMs = v
	return s
}

func (s CronJobState) WithLastRunStatus(v string) CronJobState {
	s.lastRunStatus = v
	return s
}

func (s CronJobState) WithLastError(v string) CronJobState {
	s.lastError = v
	return s
}

func (s CronJobState) WithLastDurationMs(v int64) CronJobState {
	s.lastDurationMs = v
	return s
}

func (s CronJobState) WithConsecutiveErrors(v int) CronJobState {
	s.consecutiveErrors = v
	return s
}

func (s CronJobState) WithLastFailureAlertAtMs(v *int64) CronJobState {
	s.lastFailureAlertAtMs = v
	return s
}

func (s CronJobState) WithScheduleErrorCount(v int) CronJobState {
	s.scheduleErrorCount = v
	return s
}

func (s CronJobState) WithLastDeliveryStatus(v string) CronJobState {
	s.lastDeliveryStatus = v
	return s
}

func (s CronJobState) WithLastDeliveryError(v string) CronJobState {
	s.lastDeliveryError = v
	return s
}
