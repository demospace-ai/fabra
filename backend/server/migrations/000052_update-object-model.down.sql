ALTER TABLE objects ALTER COLUMN namespace SET NOT NULL;
ALTER TABLE objects ALTER COLUMN table_name SET NOT NULL;
ALTER TABLE objects DROP COLUMN target_type;
ALTER TABLE objects DROP COLUMN sync_mode;
ALTER TABLE objects DROP COLUMN cursor_field;
ALTER TABLE objects DROP COLUMN primary_key;
ALTER TABLE objects DROP COLUMN frequency;
ALTER TABLE objects DROP COLUMN frequency_units;

ALTER TABLE syncs ALTER COLUMN sync_mode TYPE VARCHAR(64);
ALTER TABLE syncs ALTER COLUMN frequency_units TYPE varchar(64);