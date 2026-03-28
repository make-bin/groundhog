CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    channel_id VARCHAR(100) UNIQUE NOT NULL,
    channel_type INT NOT NULL,
    plugin_id VARCHAR(100),
    status INT DEFAULT 0,
    capabilities TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
