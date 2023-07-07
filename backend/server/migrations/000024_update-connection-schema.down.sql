DROP TABLE IF EXISTS sync_field_mappings;
DROP TABLE IF EXISTS sync_configurations;
DROP TABLE IF EXISTS model_fields;
DROP TABLE IF EXISTS models;
DROP TABLE IF EXISTS destinations;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS end_customers;

ALTER TABLE connections RENAME TO data_connections;
ALTER TABLE data_connections ADD COLUMN display_name VARCHAR(255) NOT NULL DEFAULT 'Connection';

CREATE TABLE IF NOT EXISTS sync_configurations(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    data_connection_id BIGINT NOT NULL REFERENCES data_connections(id),
    dataset_name VARCHAR(255) NOT NULL,
    table_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_configurations_organization_id_idx ON sync_configurations(organization_id);