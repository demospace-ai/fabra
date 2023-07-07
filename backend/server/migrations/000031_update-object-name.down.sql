ALTER TABLE objects RENAME TO models;
ALTER TABLE object_fields RENAME TO model_fields;
ALTER TABLE model_fields RENAME COLUMN object_id TO model_id;