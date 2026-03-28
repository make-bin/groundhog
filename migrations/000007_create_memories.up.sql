-- Try to create the vector extension; skip silently if unavailable
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS vector;
EXCEPTION WHEN OTHERS THEN
    RAISE NOTICE 'pgvector extension not available, skipping. Vector search will be disabled.';
END;
$$;

CREATE TABLE IF NOT EXISTS memories (
    id         BIGSERIAL PRIMARY KEY,
    memory_id  VARCHAR(100) NOT NULL,
    user_id    VARCHAR(100) NOT NULL,
    content    TEXT NOT NULL,
    embedding  TEXT,
    tags       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT memories_memory_id_key UNIQUE (memory_id)
);

CREATE INDEX IF NOT EXISTS idx_memories_user_id ON memories(user_id);
CREATE INDEX IF NOT EXISTS idx_memories_deleted_at ON memories(deleted_at);

-- tsvector full-text index (always available)
CREATE INDEX IF NOT EXISTS idx_memories_content_fts
    ON memories USING GIN (to_tsvector('simple', content));

-- Vector index created only if pgvector is available
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector') THEN
        EXECUTE 'ALTER TABLE memories ALTER COLUMN embedding TYPE VECTOR(1024) USING embedding::vector';
        EXECUTE 'CREATE INDEX IF NOT EXISTS idx_memories_embedding ON memories USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100)';
    END IF;
END;
$$;
