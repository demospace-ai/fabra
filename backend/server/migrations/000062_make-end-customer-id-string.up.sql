ALTER TABLE link_tokens ALTER COLUMN end_customer_id SET DATA TYPE VARCHAR(256);
ALTER TABLE end_customer_api_keys ALTER COLUMN end_customer_id SET DATA TYPE VARCHAR(256);
ALTER TABLE sources ALTER COLUMN end_customer_id SET DATA TYPE VARCHAR(256);
ALTER TABLE syncs ALTER COLUMN end_customer_id SET DATA TYPE VARCHAR(256);