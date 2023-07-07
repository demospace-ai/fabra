DROP TABLE IF EXISTS sync_configurations;

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

ALTER TABLE organizations ADD COLUMN default_data_connection_id BIGINT REFERENCES data_connections(id);
ALTER TABLE organizations ADD COLUMN default_event_set_id BIGINT REFERENCES event_sets(id);

CREATE TABLE IF NOT EXISTS analyses(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    analysis_type VARCHAR(128) NOT NULL,

    connection_id BIGINT REFERENCES data_connections(id),
    event_set_id BIGINT REFERENCES event_sets(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    query TEXT,
    breakdown_property_name VARCHAR(255),
    breakdown_property_type VARCHAR(128),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX analyses_user_id_idx ON analyses(user_id);
CREATE INDEX analyses_organization_id_idx ON analyses(organization_id);

CREATE TABLE IF NOT EXISTS events(
    id BIGSERIAL PRIMARY KEY,
    analysis_id BIGINT NOT NULL REFERENCES analyses(id),
    event_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX events_analysis_id_idx ON events(analysis_id);

CREATE TABLE IF NOT EXISTS event_filters(
    id BIGSERIAL PRIMARY KEY,
    analysis_id BIGINT NOT NULL REFERENCES analyses(id),
    event_id BIGINT NOT NULL REFERENCES events(id),
    property_name VARCHAR(255) NOT NULL,
    property_type VARCHAR(128) NOT NULL,
    property_value VARCHAR(255) NOT NULL,
    filter_type VARCHAR(32) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX event_filters_analysis_id_idx ON event_filters(analysis_id);
CREATE INDEX event_filters_step_id_idx ON event_filters(event_id);

CREATE TABLE IF NOT EXISTS dashboards(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX dashboards_user_id_idx ON dashboards(user_id);
CREATE INDEX dashboards_organization_id_idx ON dashboards(organization_id);

CREATE TABLE IF NOT EXISTS dashboard_panels(
    id BIGSERIAL PRIMARY KEY,
    dashboard_id BIGINT NOT NULL REFERENCES dashboards(id),
    title VARCHAR(255) NOT NULL,
    panel_type VARCHAR(32) NOT NULL,
    analysis_id BIGINT REFERENCES analyses(id),
    content TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX dashboard_panels_analysis_id_idx ON dashboard_panels(analysis_id);
CREATE INDEX dashboard_panels_dashboard_id_idx ON dashboard_panels(dashboard_id);