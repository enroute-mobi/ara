-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lines RENAME COLUMN object_ids TO codes;
ALTER TABLE operators RENAME COLUMN object_ids TO codes;
ALTER TABLE stop_areas RENAME COLUMN object_ids TO codes;
ALTER TABLE stop_visits RENAME COLUMN object_ids TO codes;
ALTER TABLE vehicle_journeys RENAME COLUMN object_ids TO codes;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lines RENAME COLUMN codes TO object_ids;
ALTER TABLE operators RENAME COLUMN codes TO object_ids;
ALTER TABLE stop_areas RENAME COLUMN codes TO object_ids;
ALTER TABLE stop_visits RENAME COLUMN codes TO object_ids;
ALTER TABLE vehicle_journeys RENAME COLUMN codes TO object_ids;
