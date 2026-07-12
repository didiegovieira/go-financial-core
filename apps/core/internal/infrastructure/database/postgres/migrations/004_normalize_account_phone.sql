UPDATE accounts
SET phone = regexp_replace(phone, '\D', '', 'g')
WHERE phone <> regexp_replace(phone, '\D', '', 'g');
