package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/make-bin/groundhog/pkg/application/dto"
	"github.com/make-bin/groundhog/pkg/application/service"
)

// CronRPCHandler handles WebSocket RPC calls for cron job management.
type CronRPCHandler struct {
	CronAppSvc service.CronAppService `inject:""`
}

// NewCronRPCHandler creates a new CronRPCHandler.
func NewCronRPCHandler() *CronRPCHandler {
	return &CronRPCHandler{}
}

// HandleCronList handles the "cron.list" RPC method.
// Parses ListCronJobsRequest and returns a paginated list of cron jobs.
func (h *CronRPCHandler) HandleCronList(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req dto.ListCronJobsRequest
	if len(params) > 0 {
		if err := json.Unmarshal(params, &req); err != nil {
			return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
		}
	}
	return h.CronAppSvc.ListJobs(ctx, &req)
}

// HandleCronAdd handles the "cron.add" RPC method.
// Parses CreateCronJobRequest and creates a new cron job.
func (h *CronRPCHandler) HandleCronAdd(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req dto.CreateCronJobRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
	}
	if req.Name == "" {
		return nil, fmt.Errorf("INVALID_REQUEST: name is required")
	}
	if req.SessionTarget == "" {
		return nil, fmt.Errorf("INVALID_REQUEST: session_target is required")
	}
	return h.CronAppSvc.CreateJob(ctx, &req)
}

// cronUpdateParams is the request struct for the "cron.update" RPC method.
type cronUpdateParams struct {
	ID    string                   `json:"id"`
	Patch dto.UpdateCronJobRequest `json:"patch"`
}

// HandleCronUpdate handles the "cron.update" RPC method.
// Parses id and UpdateCronJobRequest patch, then updates the cron job.
func (h *CronRPCHandler) HandleCronUpdate(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p cronUpdateParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
	}
	if p.ID == "" {
		return nil, fmt.Errorf("INVALID_REQUEST: id is required")
	}
	return h.CronAppSvc.UpdateJob(ctx, p.ID, &p.Patch)
}

// cronIDParams is a minimal request struct containing only an id field.
type cronIDParams struct {
	ID string `json:"id"`
}

// HandleCronRemove handles the "cron.remove" RPC method.
// Parses id and deletes the cron job, returning {"removed": bool}.
func (h *CronRPCHandler) HandleCronRemove(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p cronIDParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
	}
	if p.ID == "" {
		return nil, fmt.Errorf("INVALID_REQUEST: id is required")
	}
	err := h.CronAppSvc.DeleteJob(ctx, p.ID)
	if err != nil {
		return map[string]bool{"removed": false}, err
	}
	return map[string]bool{"removed": true}, nil
}

// cronRunParams is the request struct for the "cron.run" RPC method.
type cronRunParams struct {
	ID   string `json:"id"`
	Mode string `json:"mode"`
}

// HandleCronRun handles the "cron.run" RPC method.
// Parses id and mode (default "due"), triggers the cron job, and returns {"triggered": bool}.
func (h *CronRPCHandler) HandleCronRun(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p cronRunParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
	}
	if p.ID == "" {
		return nil, fmt.Errorf("INVALID_REQUEST: id is required")
	}
	mode := p.Mode
	if mode == "" {
		mode = "due"
	}
	err := h.CronAppSvc.TriggerJob(ctx, p.ID, mode)
	if err != nil {
		return map[string]bool{"triggered": false}, err
	}
	return map[string]bool{"triggered": true}, nil
}

// HandleCronStatus handles the "cron.status" RPC method.
// Returns the overall cron scheduler status.
func (h *CronRPCHandler) HandleCronStatus(ctx context.Context, params json.RawMessage) (interface{}, error) {
	return h.CronAppSvc.GetStatus(ctx)
}

// cronRunsParams is the request struct for the "cron.runs" RPC method.
type cronRunsParams struct {
	Scope            string   `json:"scope"`             // "job" or "all"
	JobID            string   `json:"job_id"`
	Statuses         []string `json:"statuses"`
	DeliveryStatuses []string `json:"delivery_statuses"`
	Query            string   `json:"query"`
	SortDir          string   `json:"sort_dir"`
	Offset           int      `json:"offset"`
	Limit            int      `json:"limit"`
}

// HandleCronRuns handles the "cron.runs" RPC method.
// Parses scope/jobId/filter params and returns run logs.
// When scope is "job" (or jobId is provided), calls GetRunLogs; otherwise calls GetAllRunLogs.
func (h *CronRPCHandler) HandleCronRuns(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p cronRunsParams
	if len(params) > 0 {
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, fmt.Errorf("INVALID_REQUEST: failed to parse params: %w", err)
		}
	}

	req := &dto.GetRunLogsRequest{
		JobID:            p.JobID,
		Statuses:         p.Statuses,
		DeliveryStatuses: p.DeliveryStatuses,
		Query:            p.Query,
		SortDir:          p.SortDir,
		Offset:           p.Offset,
		Limit:            p.Limit,
	}

	if p.Scope == "job" || p.JobID != "" {
		if p.JobID == "" {
			return nil, fmt.Errorf("INVALID_REQUEST: job_id is required when scope is \"job\"")
		}
		return h.CronAppSvc.GetRunLogs(ctx, req)
	}

	return h.CronAppSvc.GetAllRunLogs(ctx, req)
}
