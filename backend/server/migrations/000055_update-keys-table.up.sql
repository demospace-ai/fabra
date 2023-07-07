ALTER TABLE api_keys RENAME COLUMN api_key TO encrypted_key;

CREATE TABLE IF NOT EXISTS webhook_keys (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    encrypted_private_key VARCHAR(256) NOT NULL,
    public_key VARCHAR(256) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX webhook_keys_organization_id_idx ON webhook_keys(organization_id);