CREATE TABLE IF NOT EXISTS verifications(
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(6) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id),
    expiration TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX verifications_code_and_user_id_idx ON verifications(code, user_id);

