ALTER TABLE object_fields ADD COLUMN display_name VARCHAR(256);
ALTER TABLE object_fields ADD COLUMN description TEXT;
ALTER TABLE object_fields ADD COLUMN omit BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE objects RENAME COLUMN customer_id_column TO end_customer_id_column;
