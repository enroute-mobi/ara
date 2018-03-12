-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE vehicle_journeys
  ADD COLUMN origin_name text,
  ADD COLUMN destination_name text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE vehicle_journeys
  DROP COLUMN IF EXISTS origin_name,
  DROP COLUMN IF EXISTS destination_name;