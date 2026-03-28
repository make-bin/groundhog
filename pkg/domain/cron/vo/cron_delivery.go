package vo

import "fmt"

// DeliveryMode identifies how cron job results are delivered.
type DeliveryMode string

const (
	DeliveryModeAnnounce DeliveryMode = "announce"
	DeliveryModeWebhook  DeliveryMode = "webhook"
	DeliveryModeNone     DeliveryMode = "none"
)

// CronDelivery holds result delivery configuration for a cron job.
// It is immutable after construction.
type CronDelivery struct {
	mode               DeliveryMode
	channel            string
	to                 string
	accountId          string
	bestEffort         bool
	failureDestination string
}

// NewCronDelivery creates a CronDelivery value object.
// mode must be one of announce, webhook, or none.
func NewCronDelivery(mode DeliveryMode, channel, to, accountId string, bestEffort bool, failureDestination string) (CronDelivery, error) {
	switch mode {
	case DeliveryModeAnnounce, DeliveryModeWebhook, DeliveryModeNone:
		// valid
	default:
		return CronDelivery{}, fmt.Errorf("invalid delivery mode %q: must be announce, webhook, or none", mode)
	}
	return CronDelivery{
		mode:               mode,
		channel:            channel,
		to:                 to,
		accountId:          accountId,
		bestEffort:         bestEffort,
		failureDestination: failureDestination,
	}, nil
}

// Getters

func (d CronDelivery) Mode() DeliveryMode         { return d.mode }
func (d CronDelivery) Channel() string            { return d.channel }
func (d CronDelivery) To() string                 { return d.to }
func (d CronDelivery) AccountId() string          { return d.accountId }
func (d CronDelivery) BestEffort() bool           { return d.bestEffort }
func (d CronDelivery) FailureDestination() string { return d.failureDestination }
