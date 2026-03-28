CREATE TABLE IF NOT EXISTS cron_jobs (
    id              BIGSERIAL PRIMARY KEY,
    job_id          VARCHAR(100) NOT NULL,
    agent_id        VARCHAR(100),
    session_key     VARCHAR(200),
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    enabled         BOOLEAN NOT NULL DEFAULT TRUE,
    delete_after_run BOOLEAN,
    schedule        JSONB NOT NULL,
    session_target  VARCHAR(200) NOT NULL,
    wake_mode       VARCHAR(50) NOT NULL DEFAULT 'next-heartbeat',
    payload         JSONB NOT NULL,
    delivery        JSONB,
    failure_alert   JSONB,
    state           JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT cron_jobs_job_id_key UNIQUE (job_id),
    CONSTRAINT cron_jobs_name_key UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_cron_jobs_enabled ON cron_jobs(enabled);
CREATE INDEX IF NOT EXISTS idx_cron_jobs_agent_id ON cron_jobs(agent_id);
CREATE INDEX IF NOT EXISTS idx_cron_jobs_deleted_at ON cron_jobs(deleted_at);
CREATE INDEX IF NOT EXISTS idx_cron_jobs_enabled_next_run ON cron_jobs(enabled, (state->>'nextRunAtMs'));
