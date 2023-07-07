ALTER TABLE syncs RENAME COLUMN object_id TO model_id;
ALTER TABLE sync_field_mappings RENAME COLUMN sync_id TO sync_configuration_id;