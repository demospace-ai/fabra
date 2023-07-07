ALTER TABLE event_sets ALTER COLUMN table_name DROP NOT NULL;
ALTER TABLE event_sets ALTER COLUMN dataset_name DROP NOT NULL;
ALTER TABLE event_sets ADD CONSTRAINT check_table_configuration CHECK ((dataset_name IS NOT NULL AND table_name IS NOT NULL) or custom_join IS NOT NULL)