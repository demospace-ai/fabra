CREATE TABLE IF NOT EXISTS data_connections(
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    connection_type VARCHAR(255) NOT NULL,

    username VARCHAR(255),
    password VARCHAR(255),
    credentials VARCHAR(1024),
    warehouse_name VARCHAR(255),
    database_name VARCHAR(255),
    role VARCHAR(255),
    account VARCHAR(255),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX data_connections_organization_id_idx ON data_connections(organization_id);