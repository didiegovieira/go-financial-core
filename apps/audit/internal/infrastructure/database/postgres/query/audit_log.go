package query

const AuditLogCreate = `
INSERT INTO audit_logs (
	id,
	event_name,
	aggregate_id,
	payload,
	occurred_at
) VALUES ($1, $2, $3, $4, $5);
`
