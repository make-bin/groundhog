package vo

import "fmt"

// CronFailureAlert holds failure alert configuration for a cron job.
// It is immutable after construction.
type CronFailureAlert struct {
	after      int // default 3
	channel    string
	to         string
	cooldownMs int64
	mode       string // announce/webhook
	accountId  string
}

// NewCronFailureAlert creates a CronFailureAlert value object.
// If after <= 0 it defaults to 3. mode must be "announce" or "webhook".
func NewCronFailureAlert(after int, channel, to string, cooldownMs int64, mode, accountId string) (CronFailureAlert, error) {
	if after <= 0 {
		after = 3
	}
	switch mode {
	case "announce", "webhook":
		// valid
	default:
		return CronFailureAlert{}, fmt.Errorf("invalid failure alert mode %q: must be announce or webhook", mode)
	}
	return CronFailureAlert{
		after:      after,
		channel:    channel,
		to:         to,
		cooldownMs: cooldownMs,
		mode:       mode,
		accountId:  accountId,
	}, nil
}

// Getters

func (a CronFailureAlert) After() int        { return a.after }
func (a CronFailureAlert) Channel() string   { return a.channel }
func (a CronFailureAlert) To() string        { return a.to }
func (a CronFailureAlert) CooldownMs() int64 { return a.cooldownMs }
func (a CronFailureAlert) Mode() string      { return a.mode }
func (a CronFailureAlert) AccountId() string { return a.accountId }
