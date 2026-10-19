-- Local demo only. Reapplying preserves existing accounts and balances.
INSERT INTO accounts (
    id, name, email, password, phone,
    document_type, document_number, document_country
) VALUES
    ('10000000-0000-4000-8000-000000000001', 'Demo Alice', 'alice@example.test',
     'local-demo-only', '5511999990001', 'CPF', '00000000001', 'BR'),
    ('10000000-0000-4000-8000-000000000002', 'Demo Bruno', 'bruno@example.test',
     'local-demo-only', '5511999990002', 'CPF', '00000000002', 'BR')
ON CONFLICT DO NOTHING;

-- Balance is stored in cents: 100000 = 1000.00, 25000 = 250.00.
INSERT INTO wallets (id, account_id, balance, status)
SELECT seed.id, seed.account_id, seed.balance, 'active'
FROM (VALUES
    ('20000000-0000-4000-8000-000000000001'::uuid,
     '10000000-0000-4000-8000-000000000001'::uuid, 100000::bigint),
    ('20000000-0000-4000-8000-000000000002'::uuid,
     '10000000-0000-4000-8000-000000000002'::uuid, 25000::bigint)
) AS seed(id, account_id, balance)
JOIN accounts ON accounts.id = seed.account_id
WHERE NOT EXISTS (SELECT 1 FROM wallets WHERE wallets.account_id = seed.account_id)
ON CONFLICT DO NOTHING;
