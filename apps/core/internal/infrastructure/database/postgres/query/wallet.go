package query

const WalletCreate = `
INSERT INTO wallets (
	id,
	account_id,
	balance,
	status
) VALUES ($1, $2, $3, $4);
`
