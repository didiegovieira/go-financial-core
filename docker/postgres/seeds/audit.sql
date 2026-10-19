-- Synthetic history for local exploration, not evidence of Kafka delivery.
INSERT INTO audit_logs (id, event_name, aggregate_id, payload, occurred_at)
VALUES (
    '30000000-0000-4000-8000-000000000001',
    'account.created',
    '10000000-0000-4000-8000-000000000001',
    '{"source":"local-demo-seed","account_id":"10000000-0000-4000-8000-000000000001","name":"Demo Alice"}'::jsonb,
    '2026-01-01 12:00:00'
)
ON CONFLICT DO NOTHING;
