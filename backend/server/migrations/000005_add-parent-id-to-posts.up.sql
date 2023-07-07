ALTER TABLE posts ADD parent_post_id BIGINT REFERENCES posts(id);
CREATE INDEX parent_post_id_idx ON posts(parent_post_id);