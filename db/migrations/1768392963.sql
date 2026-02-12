-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE lines RENAME COLUMN attributes TO raw_attributes;
ALTER TABLE stop_areas RENAME COLUMN attributes TO raw_attributes;
ALTER TABLE stop_visits RENAME COLUMN attributes TO raw_attributes;
ALTER TABLE vehicle_journeys RENAME COLUMN attributes TO raw_attributes;


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE lines RENAME COLUMN raw_attributes TO attributes;
ALTER TABLE stop_areas RENAME COLUMN raw_attributes TO attributes;
ALTER TABLE stop_visits RENAME COLUMN raw_attributes TO attributes;
ALTER TABLE vehicle_journeys RENAME COLUMN raw_attributes TO attributes;
