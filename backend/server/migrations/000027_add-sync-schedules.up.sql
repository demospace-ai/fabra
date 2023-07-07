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