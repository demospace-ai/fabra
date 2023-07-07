CREATE TABLE IF NOT EXISTS posts(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    post_type VARCHAR NOT NULL,
    body text NOT NULL,
    title text,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deactivated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX post_user_id_idx ON posts(user_id);