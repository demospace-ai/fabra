ALTER TABLE syncs RENAME COLUMN model_id TO object_id;
ALTER TABLE sync_field_mappings RENAME COLUMN sync_configuration_id TO sync_id;