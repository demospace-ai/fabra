CREATE TABLE IF NOT EXISTS organizations(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email_domain VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE users ADD organization_id BIGINT REFERENCES organizations(id);
ALTER TABLE users ADD email VARCHAR(255) NOT NULL;
ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;
ALTER TABLE users ALTER COLUMN last_name SET NOT NULL;
ALTER TABLE posts ADD organization_id BIGINT REFERENCES organizations(id) NOT NULL;

DROP TABLE emails;

CREATE INDEX users_organization_id_idx ON users(organization_id);
CREATE INDEX posts_organization_id_idx ON posts(organization_id);
CREATE INDEX organizations_email_domain_idx ON organizations(email_domain);