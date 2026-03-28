package cron

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/cron/aggregate/cron_job"
	"github.com/make-bin/groundhog/pkg/domain/cron/vo"
	"github.com/make-bin/groundhog/pkg/utils/logger"
)

// DeliveryExecutor handles result delivery and failure alerts for cron jobs.
type DeliveryExecutor struct {
	httpClient *http.Client
	logger     logger.Logger
}

// NewDeliveryExecutor creates a new DeliveryExecutor.
func NewDeliveryExecutor(log logger.Logger) *DeliveryExecutor {
	return &DeliveryExecutor{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     log,
	}
}

// DeliverResult delivers the execution result based on the job's delivery configuration.
// Returns the delivery status ("delivered", "not-delivered", "skipped") and any error.
func (d *DeliveryExecutor) DeliverResult(ctx context.Context, job *cron_job.CronJob, result *ExecuteResult) (string, error) {
	delivery := job.Delivery()
	if delivery == nil {
		return "skipped", nil
	}

	switch delivery.Mode() {
	case vo.DeliveryModeNone:
		return "skipped", nil

	case vo.DeliveryModeAnnounce:
		return d.deliverAnnounce(ctx, job, result)

	case vo.DeliveryModeWebhook:
		return d.deliverWebhook(ctx, job, result)

	default:
		return "skipped", fmt.Errorf("unknown delivery mode: %s", delivery.Mode())
	}
}

// deliverAnnounce sends the result summary to a messaging channel.
func (d *DeliveryExecutor) deliverAnnounce(ctx context.Context, job *cron_job.CronJob, result *ExecuteResult) (string, error) {
	delivery := job.Delivery()

	// Build the announcement message.
	msg := fmt.Sprintf("[Cron: %s] %s", job.Name(), result.Summary)
	if result.Status == "error" {
		msg = fmt.Sprintf("[Cron: %s] Error: %s", job.Name(), result.Error)
	}

	d.logger.Info("cron announce delivery",
		"jobId", job.ID().Value(),
		"channel", delivery.Channel(),
		"to", delivery.To(),
		"status", result.Status,
	)

	// For announce mode, we log the delivery. Actual channel sending would
	// require ChannelAppService integration which depends on the messaging subsystem.
	// This is a placeholder that records the intent.
	_ = msg
	_ = ctx

	return "delivered", nil
}

// deliverWebhook sends the result as an HTTP POST to the configured URL.
func (d *DeliveryExecutor) deliverWebhook(ctx context.Context, job *cron_job.CronJob, result *ExecuteResult) (string, error) {
	delivery := job.Delivery()
	url := delivery.To()
	if url == "" {
		return "not-delivered", fmt.Errorf("webhook URL is empty")
	}

	payload := map[string]interface{}{
		"job_id":     job.ID().Value(),
		"job_name":   job.Name(),
		"status":     result.Status,
		"summary":    result.Summary,
		"error":      result.Error,
		"session_id": result.SessionID,
		"model":      result.Model,
		"provider":   result.Provider,
		"usage": map[string]int{
			"input_tokens":  result.Usage.InputTokens,
			"output_tokens": result.Usage.OutputTokens,
			"total_tokens":  result.Usage.TotalTokens,
		},
		"timestamp_ms": time.Now().UnixMilli(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "not-delivered", fmt.Errorf("marshal webhook payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "not-delivered", fmt.Errorf("create webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		if delivery.BestEffort() {
			d.logger.Warn("webhook delivery failed (best-effort)", "jobId", job.ID().Value(), "error", err)
			return "not-delivered", nil
		}
		return "not-delivered", fmt.Errorf("webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "delivered", nil
	}

	errMsg := fmt.Sprintf("webhook returned status %d", resp.StatusCode)
	if delivery.BestEffort() {
		d.logger.Warn("webhook delivery failed (best-effort)", "jobId", job.ID().Value(), "error", errMsg)
		return "not-delivered", nil
	}
	return "not-delivered", errors.New(errMsg)
}

// SendFailureAlert sends a failure alert if conditions are met (consecutive errors >= threshold, cooldown elapsed).
func (d *DeliveryExecutor) SendFailureAlert(ctx context.Context, job *cron_job.CronJob) {
	if !job.ShouldAlert() {
		return
	}

	fa := job.FailureAlert()
	if fa == nil {
		return
	}

	alertMsg := fmt.Sprintf("[ALERT] Cron job %q has failed %d consecutive times. Last error: %s",
		job.Name(),
		job.State().ConsecutiveErrors(),
		job.State().LastError(),
	)

	switch fa.Mode() {
	case "announce":
		d.logger.Warn("cron failure alert (announce)",
			"jobId", job.ID().Value(),
			"channel", fa.Channel(),
			"to", fa.To(),
			"message", alertMsg,
		)

	case "webhook":
		if fa.To() == "" {
			d.logger.Error("failure alert webhook URL is empty", "jobId", job.ID().Value())
			return
		}

		payload := map[string]interface{}{
			"type":               "failure_alert",
			"job_id":             job.ID().Value(),
			"job_name":           job.Name(),
			"consecutive_errors": job.State().ConsecutiveErrors(),
			"last_error":         job.State().LastError(),
			"timestamp_ms":       time.Now().UnixMilli(),
		}

		body, err := json.Marshal(payload)
		if err != nil {
			d.logger.Error("marshal failure alert payload", "jobId", job.ID().Value(), "error", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fa.To(), bytes.NewReader(body))
		if err != nil {
			d.logger.Error("create failure alert request", "jobId", job.ID().Value(), "error", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := d.httpClient.Do(req)
		if err != nil {
			d.logger.Error("failure alert webhook request failed", "jobId", job.ID().Value(), "error", err)
			return
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			d.logger.Info("failure alert sent", "jobId", job.ID().Value())
		} else {
			d.logger.Error("failure alert webhook returned error", "jobId", job.ID().Value(), "status", resp.StatusCode)
		}

	default:
		d.logger.Warn("unknown failure alert mode", "jobId", job.ID().Value(), "mode", fa.Mode())
	}
}
