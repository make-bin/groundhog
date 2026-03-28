CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(100) UNIQUE NOT NULL,
    agent_id VARCHAR(100) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    active_model TEXT,
    state INT DEFAULT 0,
    metadata TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_deleted_at ON sessions(deleted_at);

CREATE TABLE turns (
    id SERIAL PRIMARY KEY,
    session_po_id INT REFERENCES sessions(id) ON DELETE CASCADE,
    turn_id VARCHAR(100),
    user_input TEXT,
    response TEXT,
    model_used VARCHAR(100),
    token_usage TEXT,
    tool_calls TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP
);
CREATE INDEX idx_turns_session_po_id ON turns(session_po_id);
