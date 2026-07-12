CREATE TABLE IF NOT EXISTS audit_logs (
	id UUID PRIMARY KEY,
	event_name VARCHAR(255) NOT NULL,
	aggregate_id VARCHAR(255) NOT NULL,
	payload JSONB NOT NULL,
	occurred_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_event_name ON audit_logs(event_name);
CREATE INDEX IF NOT EXISTS idx_audit_logs_aggregate_id ON audit_logs(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_occurred_at ON audit_logs(occurred_at);
