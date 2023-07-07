ALTER TABLE event_sets ALTER COLUMN table_name SET NOT NULL;
ALTER TABLE event_sets ALTER COLUMN dataset_name SET NOT NULL;
ALTER TABLE event_sets DROP CONSTRAINT check_table_configuration;