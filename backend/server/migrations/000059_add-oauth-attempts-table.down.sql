ALTER TABLE external_profiles DROP COLUMN oauth_provider;
ALTER TABLE users ADD COLUMN first_name VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN last_name VARCHAR(256) NOT NULL DEFAULT '';
UPDATE users SET first_name = split_part(name, ' ', 1);
UPDATE users SET last_name = split_part(name, ' ', 2);
ALTER TABLE users ALTER COLUMN first_name DROP DEFAULT;
ALTER TABLE users ALTER COLUMN last_name DROP DEFAULT;
ALTER TABLE users DROP COLUMN name;