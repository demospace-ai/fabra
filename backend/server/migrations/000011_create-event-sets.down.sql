DROP TABLE IF EXISTS event_sets;
ALTER TABLE data_connections ALTER COLUMN organization_id TYPE BIGINT NOT NULL;