CREATE TABLE plugins (
    id SERIAL PRIMARY KEY,
    plugin_id VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    version VARCHAR(50),
    plugin_type VARCHAR(50),
    entry_point VARCHAR(500),
    status INT DEFAULT 0,
    capabilities TEXT,
    started_at TIMESTAMP,
    restart_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
