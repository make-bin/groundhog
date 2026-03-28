CREATE TABLE media_assets (
    id SERIAL PRIMARY KEY,
    asset_id VARCHAR(100) UNIQUE NOT NULL,
    mime_type VARCHAR(100),
    size BIGINT,
    source INT,
    url TEXT,
    local_path TEXT,
    validated BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE INDEX idx_media_assets_asset_id ON media_assets(asset_id);
