CREATE TABLE sync_runs (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    sync_id         BIGINT NOT NULL REFERENCES syncs(id),
    status          VARCHAR(32) NOT NULL,
    error           TEXT,
    started_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_at    TIMESTAMP WITH TIME ZONE,

    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX sync_runs_sync_id_idx ON sync_runs(sync_id);