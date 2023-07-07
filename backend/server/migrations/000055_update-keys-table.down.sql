ALTER TABLE api_keys RENAME COLUMN encrypted_key TO api_key;
DROP TABLE IF EXISTS webhook_keys;