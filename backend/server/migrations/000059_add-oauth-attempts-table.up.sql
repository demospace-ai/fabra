ALTER TABLE external_profiles ADD COLUMN oauth_provider VARCHAR(64) NOT NULL DEFAULT 'google';
ALTER TABLE external_profiles ALTER COLUMN oauth_provider DROP DEFAULT;
ALTER TABLE users ADD COLUMN name VARCHAR(256) NOT NULL DEFAULT '';
UPDATE users SET name = first_name  || ' ' || last_name;
ALTER TABLE users ALTER COLUMN name DROP DEFAULT;
ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;