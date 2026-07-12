DO $$
BEGIN
	IF (
		SELECT data_type = 'numeric'
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'wallets'
			AND column_name = 'balance'
	) THEN
		ALTER TABLE wallets
		ALTER COLUMN balance TYPE BIGINT
		USING ROUND(balance * 100)::BIGINT;
	END IF;
END $$;
