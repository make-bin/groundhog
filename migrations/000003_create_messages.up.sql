CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    message_id VARCHAR(100) UNIQUE NOT NULL,
    channel_id VARCHAR(100) NOT NULL,
    account_id VARCHAR(100) NOT NULL,
    content TEXT,
    status INT DEFAULT 0,
    routed_to VARCHAR(100),
    received_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_messages_channel_id ON messages(channel_id);
CREATE INDEX idx_messages_account_id ON messages(account_id);
