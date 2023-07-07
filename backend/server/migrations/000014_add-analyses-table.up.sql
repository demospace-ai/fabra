CREATE TABLE IF NOT EXISTS analyses(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    analysis_type VARCHAR(128) NOT NULL,

    connection_id BIGINT REFERENCES data_connections(id),
    event_set_id BIGINT REFERENCES event_sets(id),
    title VARCHAR(255),
    query TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX analyses_user_id_idx ON analyses(user_id);
CREATE INDEX analyses_organization_id_idx ON analyses(organization_id);

CREATE TABLE IF NOT EXISTS funnel_steps(
    id BIGSERIAL PRIMARY KEY,
    analysis_id BIGINT NOT NULL REFERENCES analyses(id),
    step_name VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX funnel_steps_analysis_id_idx ON funnel_steps(analysis_id);