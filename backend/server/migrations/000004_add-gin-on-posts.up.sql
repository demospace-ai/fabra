CREATE INDEX fts_idx ON posts 
    USING gin((setweight(to_tsvector('english', coalesce(title, '')), 'A') || setweight(to_tsvector('english', body), 'B')))