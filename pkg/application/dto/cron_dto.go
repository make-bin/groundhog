package dto

// --- Request DTOs ---

type CronScheduleDTO struct {
	Kind      string `json:"kind"`
	At        string `json:"at,omitempty"`
	EveryMs   int64  `json:"every_ms,omitempty"`
	AnchorMs  *int64 `json:"anchor_ms,omitempty"`
	Expr      string `json:"expr,omitempty"`
	Tz        string `json:"tz,omitempty"`
	StaggerMs int64  `json:"stagger_ms,omitempty"`
}

type CronPayloadDTO struct {
	Kind           string `json:"kind"`
	Text           string `json:"text,omitempty"`
	Message        string `json:"message,omitempty"`
	Model          string `json:"model,omitempty"`
	Thinking       bool   `json:"thinking,omitempty"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
	LightContext   bool   `json:"light_context,omitempty"`
}

type CronDeliveryDTO struct {
	Mode               string `json:"mode"`
	Channel            string `json:"channel,omitempty"`
	To                 string `json:"to,omitempty"`
	AccountId          string `json:"account_id,omitempty"`
	BestEffort         bool   `json:"best_effort,omitempty"`
	FailureDestination string `json:"failure_destination,omitempty"`
}

type CronFailureAlertDTO struct {
	After      int    `json:"after,omitempty"`
	Channel    string `json:"channel,omitempty"`
	To         string `json:"to,omitempty"`
	CooldownMs int64  `json:"cooldown_ms,omitempty"`
	Mode       string `json:"mode,omitempty"`
	AccountId  string `json:"account_id,omitempty"`
}

type CreateCronJobRequest struct {
	AgentId        string               `json:"agent_id,omitempty"`
	SessionKey     string               `json:"session_key,omitempty"`
	Name           string               `json:"name" binding:"required"`
	Description    string               `json:"description,omitempty"`
	Enabled        *bool                `json:"enabled,omitempty"`
	DeleteAfterRun *bool                `json:"delete_after_run,omitempty"`
	Schedule       CronScheduleDTO      `json:"schedule" binding:"required"`
	SessionTarget  string               `json:"session_target" binding:"required"`
	WakeMode       string               `json:"wake_mode,omitempty"`
	Payload        CronPayloadDTO       `json:"payload" binding:"required"`
	Delivery       *CronDeliveryDTO     `json:"delivery,omitempty"`
	FailureAlert   *CronFailureAlertDTO `json:"failure_alert,omitempty"`
}

type UpdateCronJobRequest struct {
	Name           *string              `json:"name,omitempty"`
	Description    *string              `json:"description,omitempty"`
	Enabled        *bool                `json:"enabled,omitempty"`
	DeleteAfterRun *bool                `json:"delete_after_run,omitempty"`
	Schedule       *CronScheduleDTO     `json:"schedule,omitempty"`
	SessionTarget  *string              `json:"session_target,omitempty"`
	WakeMode       *string              `json:"wake_mode,omitempty"`
	Payload        *CronPayloadDTO      `json:"payload,omitempty"`
	Delivery       *CronDeliveryDTO     `json:"delivery,omitempty"`
	FailureAlert   *CronFailureAlertDTO `json:"failure_alert,omitempty"`
}

type ListCronJobsRequest struct {
	IncludeDisabled bool   `json:"include_disabled" form:"include_disabled"`
	Enabled         *bool  `json:"enabled" form:"enabled"`
	Query           string `json:"query" form:"query"`
	SortBy          string `json:"sort_by" form:"sort_by"`
	SortDir         string `json:"sort_dir" form:"sort_dir"`
	Offset          int    `json:"offset" form:"offset"`
	Limit           int    `json:"limit" form:"limit"`
}

type GetRunLogsRequest struct {
	JobID            string   `json:"job_id" form:"job_id"`
	Statuses         []string `json:"statuses" form:"statuses"`
	DeliveryStatuses []string `json:"delivery_statuses" form:"delivery_statuses"`
	Query            string   `json:"query" form:"query"`
	SortDir          string   `json:"sort_dir" form:"sort_dir"`
	Offset           int      `json:"offset" form:"offset"`
	Limit            int      `json:"limit" form:"limit"`
}

// --- Response DTOs ---

type CronJobStateDTO struct {
	NextRunAtMs          *int64 `json:"next_run_at_ms"`
	RunningAtMs          *int64 `json:"running_at_ms"`
	LastRunAtMs          *int64 `json:"last_run_at_ms"`
	LastRunStatus        string `json:"last_run_status"`
	LastError            string `json:"last_error"`
	LastDurationMs       int64  `json:"last_duration_ms"`
	ConsecutiveErrors    int    `json:"consecutive_errors"`
	LastFailureAlertAtMs *int64 `json:"last_failure_alert_at_ms"`
	ScheduleErrorCount   int    `json:"schedule_error_count"`
	LastDeliveryStatus   string `json:"last_delivery_status"`
	LastDeliveryError    string `json:"last_delivery_error"`
}

type CronJobResponse struct {
	ID             string               `json:"id"`
	AgentId        string               `json:"agent_id,omitempty"`
	SessionKey     string               `json:"session_key,omitempty"`
	Name           string               `json:"name"`
	Description    string               `json:"description,omitempty"`
	Enabled        bool                 `json:"enabled"`
	DeleteAfterRun *bool                `json:"delete_after_run,omitempty"`
	CreatedAtMs    int64                `json:"created_at_ms"`
	UpdatedAtMs    int64                `json:"updated_at_ms"`
	Schedule       CronScheduleDTO      `json:"schedule"`
	SessionTarget  string               `json:"session_target"`
	WakeMode       string               `json:"wake_mode"`
	Payload        CronPayloadDTO       `json:"payload"`
	Delivery       *CronDeliveryDTO     `json:"delivery,omitempty"`
	FailureAlert   *CronFailureAlertDTO `json:"failure_alert,omitempty"`
	State          CronJobStateDTO      `json:"state"`
}

type CronJobListResponse struct {
	Jobs   []*CronJobResponse `json:"jobs"`
	Total  int                `json:"total"`
	Offset int                `json:"offset"`
	Limit  int                `json:"limit"`
}

type CronRunLogResponse struct {
	ID             int64  `json:"id"`
	JobID          string `json:"job_id"`
	Ts             int64  `json:"ts"`
	Action         string `json:"action"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Summary        string `json:"summary,omitempty"`
	SessionID      string `json:"session_id,omitempty"`
	RunAtMs        int64  `json:"run_at_ms"`
	DurationMs     int64  `json:"duration_ms"`
	NextRunAtMs    *int64 `json:"next_run_at_ms,omitempty"`
	Model          string `json:"model,omitempty"`
	Provider       string `json:"provider,omitempty"`
	InputTokens    int    `json:"input_tokens,omitempty"`
	OutputTokens   int    `json:"output_tokens,omitempty"`
	TotalTokens    int    `json:"total_tokens,omitempty"`
	DeliveryStatus string `json:"delivery_status,omitempty"`
	DeliveryError  string `json:"delivery_error,omitempty"`
}

type RunLogListResponse struct {
	Logs   []*CronRunLogResponse `json:"logs"`
	Total  int                   `json:"total"`
	Offset int                   `json:"offset"`
	Limit  int                   `json:"limit"`
}

type CronSchedulerStatusResponse struct {
	Running     bool   `json:"running"`
	EnabledJobs int    `json:"enabled_jobs"`
	RunningJobs int    `json:"running_jobs"`
	NextRunAtMs *int64 `json:"next_run_at_ms"`
	HeartbeatMs int64  `json:"heartbeat_ms"`
}
