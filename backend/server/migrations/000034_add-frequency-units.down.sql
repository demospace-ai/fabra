ALTER TABLE syncs DROP COLUMN frequency_units;
ALTER TABLE syncs DROP CONSTRAINT check_source_configuration;
ALTER TABLE syncs ALTER COLUMN frequency DROP NOT NULL;
ALTER TABLE syncs ALTER COLUMN sync_mode DROP NOT NULL;
ALTER TABLE syncs DROP COLUMN cursor_position;