DROP TABLE IF EXISTS sync_schedules;

ALTER TABLE sync_configurations RENAME TO syncs;
ALTER TABLE syncs ADD COLUMN namespace VARCHAR(255);
ALTER TABLE syncs ADD COLUMN table_name VARCHAR(255);
ALTER TABLE syncs ADD COLUMN custom_join TEXT;
ALTER TABLE syncs ADD COLUMN frequency INT;
ALTER TABLE syncs ADD COLUMN sync_mode VARCHAR(64);
ALTER TABLE syncs ADD COLUMN cursor_field VARCHAR(255);
ALTER TABLE syncs ADD COLUMN primary_key VARCHAR(255);

ALTER TABLE sources DROP COLUMN namespace;
ALTER TABLE sources DROP COLUMN table_name;
ALTER TABLE sources DROP COLUMN custom_join;

ALTER TABLE syncs DROP COLUMN end_customer_id;