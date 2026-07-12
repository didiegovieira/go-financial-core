package query

const WalletCreate = `
INSERT INTO wallets (
	id,
	account_id,
	balance,
	status
) VALUES ($1, $2, $3, $4);
`

const WalletDeleteByAccountID = `
DELETE FROM wallets
WHERE account_id = $1;
`

const WalletCredit = `
UPDATE wallets
SET balance = balance + $1,
	updated_at = NOW()
WHERE account_id = $2
	AND status = 'active'
RETURNING balance;
`

const WalletDebit = `
UPDATE wallets
SET balance = balance - $1,
	updated_at = NOW()
WHERE account_id = $2
	AND status = 'active'
	AND balance >= $1
RETURNING balance;
`
