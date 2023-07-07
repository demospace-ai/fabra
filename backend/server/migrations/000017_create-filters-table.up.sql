CREATE TABLE IF NOT EXISTS step_filters(
    id BIGSERIAL PRIMARY KEY,
    analysis_id BIGINT NOT NULL REFERENCES analyses(id),
    step_id BIGINT NOT NULL REFERENCES funnel_steps(id),
    property_name VARCHAR(255) NOT NULL,
    property_type VARCHAR(128) NOT NULL,
    property_value VARCHAR(255) NOT NULL,
    filter_type VARCHAR(32) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX step_filters_analysis_id_idx ON step_filters(analysis_id);
CREATE INDEX step_filters_step_id_idx ON step_filters(step_id);