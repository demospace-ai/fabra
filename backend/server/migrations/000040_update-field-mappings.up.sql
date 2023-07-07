ALTER TABLE sync_field_mappings RENAME COLUMN destination_field_name TO destination_field_id;
ALTER TABLE sync_field_mappings ALTER COLUMN destination_field_id TYPE BIGINT USING destination_field_id::BIGINT;
ALTER TABLE sync_field_mappings ADD CONSTRAINT fkey_destination_field_id_object_fields_id FOREIGN KEY (destination_field_id) REFERENCES object_fields(id);