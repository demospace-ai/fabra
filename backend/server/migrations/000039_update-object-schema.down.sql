ALTER TABLE object_fields DROP COLUMN display_name;
ALTER TABLE object_fields DROP COLUMN description;
ALTER TABLE object_fields DROP COLUMN omit;
ALTER TABLE objects RENAME COLUMN end_customer_id_column TO customer_id_column;