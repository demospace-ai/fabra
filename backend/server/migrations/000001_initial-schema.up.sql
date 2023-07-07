CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS sessions(
    id BIGSERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    expiration TIMESTAMP WITH TIME ZONE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS external_profiles(
    id BIGSERIAL PRIMARY KEY,
    external_id VARCHAR(255),
    user_id BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS user_identities(
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    user_id BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS emails(
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    user_id BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX user_identities_user_id_idx ON user_identities(user_id);
CREATE INDEX emails_user_id_idx ON emails(user_id);
CREATE INDEX sessions_token_idx ON sessions(token);
CREATE INDEX external_profiles_external_id_idx ON external_profiles(external_id);
