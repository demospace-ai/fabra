DROP TABLE IF EXISTS dashboard_panels;
DROP TABLE IF EXISTS dashboards;
DROP TABLE IF EXISTS event_filters;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS analyses;

CREATE TABLE IF NOT EXISTS sync_configurations(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    connection_id BIGINT NOT NULL REFERENCES data_connections(id),
    dataset_name VARCHAR(255) NOT NULL,
    table_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_configurations_organization_id_idx ON sync_configurations(organization_id);

ALTER TABLE organizations DROP COLUMN default_data_connection_id;
ALTER TABLE organizations DROP COLUMN default_event_set_id;
DROP TABLE IF EXISTS event_sets;