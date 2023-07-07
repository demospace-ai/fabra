ALTER TABLE link_tokens ALTER COLUMN end_customer_id SET DATA TYPE BIGINT USING end_customer_id::BIGINT;
ALTER TABLE end_customer_api_keys ALTER COLUMN end_customer_id SET DATA TYPE BIGINT USING end_customer_id::BIGINT;
ALTER TABLE sources ALTER COLUMN end_customer_id SET DATA TYPE BIGINT USING end_customer_id::BIGINT;
ALTER TABLE syncs ALTER COLUMN end_customer_id SET DATA TYPE BIGINT USING end_customer_id::BIGINT;