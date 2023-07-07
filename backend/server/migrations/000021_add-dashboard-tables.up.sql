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