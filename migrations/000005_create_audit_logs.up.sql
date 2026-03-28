CREATE TABLE audit_logs (
    id SERIAL PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    principal_id VARCHAR(100),
    resource_type VARCHAR(100),
    resource_id VARCHAR(100),
    details TEXT,
    source_ip VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_principal_id ON audit_logs(principal_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
