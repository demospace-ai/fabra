ALTER TABLE models RENAME TO objects;
ALTER TABLE model_fields RENAME TO object_fields;
ALTER TABLE object_fields RENAME COLUMN model_id TO object_id;