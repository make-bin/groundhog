CREATE TABLE IF NOT EXISTS cron_run_logs (
    id              BIGSERIAL PRIMARY KEY,
    job_id          VARCHAR(100) NOT NULL,
    ts              BIGINT NOT NULL,
    action          VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL,
    error           TEXT,
    summary         TEXT,
    session_id      VARCHAR(200),
    run_at_ms       BIGINT NOT NULL,
    duration_ms     BIGINT,
    next_run_at_ms  BIGINT,
    model           VARCHAR(100),
    provider        VARCHAR(100),
    input_tokens    INTEGER,
    output_tokens   INTEGER,
    total_tokens    INTEGER,
    delivery_status VARCHAR(50),
    delivery_error  TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cron_run_logs_job_id ON cron_run_logs(job_id);
CREATE INDEX IF NOT EXISTS idx_cron_run_logs_status ON cron_run_logs(status);
CREATE INDEX IF NOT EXISTS idx_cron_run_logs_run_at_ms ON cron_run_logs(run_at_ms DESC);
