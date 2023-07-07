CREATE TABLE IF NOT EXISTS end_customer_api_keys (
    id              BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    end_customer_id BIGINT NOT NULL,
    encrypted_key   VARCHAR(256) NOT NULL,

    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
)