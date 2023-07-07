ALTER TABLE data_connections DROP COLUMN display_name;
ALTER TABLE data_connections RENAME TO connections;

CREATE TABLE IF NOT EXISTS destinations (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    connection_id BIGINT NOT NULL REFERENCES connections(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX destinations_organization_id_idx ON destinations(organization_id);

CREATE TABLE IF NOT EXISTS models (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    display_name VARCHAR(255) NOT NULL,
    destination_id BIGINT NOT NULL REFERENCES destinations(id),
    namespace VARCHAR(255),
    table_name VARCHAR(255),
    custom_join TEXT,
    customer_id_column VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX models_organization_id_idx ON models(organization_id);

CREATE TABLE IF NOT EXISTS model_fields (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT NOT NULL REFERENCES models(id),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(128) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX model_fields_model_id_idx ON model_fields(model_id);

CREATE TABLE IF NOT EXISTS end_customers (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX end_customers_organization_id_idx ON end_customers(organization_id);

CREATE TABLE IF NOT EXISTS sources (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    end_customer_id BIGINT NOT NULL REFERENCES end_customers(id),
    display_name VARCHAR(255) NOT NULL,
    connection_id BIGINT NOT NULL REFERENCES connections(id),
    namespace VARCHAR(255),
    table_name VARCHAR(255),
    custom_join TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sources_organization_id_idx ON sources(organization_id);
CREATE INDEX sources_end_customer_id_idx ON sources(end_customer_id);

DROP TABLE IF EXISTS sync_configurations;
CREATE TABLE IF NOT EXISTS sync_configurations (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    end_customer_id BIGINT REFERENCES end_customers(id),
    display_name VARCHAR(255) NOT NULL,
    destination_id BIGINT REFERENCES destinations(id),
    source_id BIGINT REFERENCES sources(id),
    model_id BIGINT REFERENCES models(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_configurations_organization_id_idx ON sync_configurations(organization_id);
CREATE INDEX sync_configurations_end_customer_id_idx ON sync_configurations(end_customer_id);

CREATE TABLE IF NOT EXISTS sync_field_mappings (
    id BIGSERIAL PRIMARY KEY,
    sync_configuration_id BIGINT NOT NULL REFERENCES sync_configurations(id),
    source_field_name VARCHAR(255) NOT NULL,
    destination_field_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_field_mappings_sync_configuration_id_idx ON sync_field_mappings(sync_configuration_id);