ALTER TABLE syncs ADD COLUMN frequency_units VARCHAR(64) NOT NULL;
ALTER TABLE syncs ADD CONSTRAINT check_source_configuration CHECK ((namespace IS NOT NULL AND table_name IS NOT NULL) or custom_join IS NOT NULL);
ALTER TABLE syncs ALTER COLUMN frequency SET NOT NULL;
ALTER TABLE syncs ALTER COLUMN sync_mode SET NOT NULL;
ALTER TABLE syncs ADD COLUMN cursor_position VARCHAR(255);