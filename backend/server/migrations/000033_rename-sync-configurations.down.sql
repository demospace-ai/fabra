ALTER TABLE syncs DROP COLUMN namespace;
ALTER TABLE syncs DROP COLUMN table_name;
ALTER TABLE syncs DROP COLUMN custom_join;
ALTER TABLE syncs DROP COLUMN frequency;
ALTER TABLE syncs DROP COLUMN sync_mode;
ALTER TABLE syncs DROP COLUMN cursor_field;
ALTER TABLE syncs DROP COLUMN primary_key;
ALTER TABLE syncs ADD COLUMN end_customer_id BIGINT NOT NULL;

ALTER TABLE syncs RENAME TO sync_configurations;

ALTER TABLE sources ADD COLUMN namespace VARCHAR(255);
ALTER TABLE sources ADD COLUMN table_name VARCHAR(255);
ALTER TABLE sources ADD COLUMN custom_join TEXT;

CREATE TABLE IF NOT EXISTS sync_schedules (
    id BIGSERIAL PRIMARY KEY,
    sync_configuration_id BIGINT NOT NULL REFERENCES sync_configurations(id),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    schedule_type VARCHAR(256) NOT NULL,
    interval BIGINT,
    interval_units VARCHAR(64),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_schedule_sync_configuration_id_idx ON sync_schedules(sync_configuration_id);