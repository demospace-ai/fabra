CREATE TABLE IF NOT EXISTS event_sets(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    connection_id BIGINT NOT NULL REFERENCES data_connections(id),
    dataset_name VARCHAR(255) NOT NULL,
    table_name VARCHAR(255) NOT NULL,
    event_type_column VARCHAR(255) NOT NULL,
    timestamp_column VARCHAR(255) NOT NULL,
    user_identifier_column VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX event_sets_organization_id_idx ON event_sets(organization_id);

ALTER TABLE data_connections ADD FOREIGN KEY (organization_id) REFERENCES organizations(id);