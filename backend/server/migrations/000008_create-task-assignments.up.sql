ALTER TABLE posts ADD assigned_user_id BIGINT REFERENCES users(id);
CREATE INDEX posts_assigned_user_id_idx ON posts(assigned_user_id);