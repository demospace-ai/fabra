ALTER TABLE objects ALTER COLUMN frequency SET NOT NULL;
ALTER TABLE objects ALTER COLUMN frequency_units SET NOT NULL;
ALTER TABLE objects DROP COLUMN recurring;
ALTER TABLE syncs ALTER COLUMN frequency SET NOT NULL;
ALTER TABLE syncs ALTER COLUMN frequency_units SET NOT NULL;
ALTER TABLE syncs DROP COLUMN recurring;