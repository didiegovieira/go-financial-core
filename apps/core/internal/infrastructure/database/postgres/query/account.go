package query

const AccountCreate = `
INSERT INTO accounts (
	id,
	name,
	email,
	password,
	phone,
	document_type,
	document_number,
	document_country
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`

const AccountFindByDocument = `
SELECT
	id,
	name,
	email,
	password,
	phone,
	document_type,
	document_number,
	document_country,
	created_at,
	updated_at
FROM accounts
WHERE document_number = $1
	AND document_type = $2
	AND document_country = $3
LIMIT 1;
`
