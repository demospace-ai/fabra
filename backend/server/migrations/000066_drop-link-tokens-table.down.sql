CREATE TABLE IF NOT EXISTS link_tokens (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    end_customer_id BIGINT NOT NULL,
    hashed_token VARCHAR(64) NOT NULL DEFAULT '',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX link_tokens_hashed_token_idx ON link_tokens(hashed_token);