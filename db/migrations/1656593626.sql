-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE vehicle_journeys
      ADD COLUMN direction_type text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE vehicle_journeys
      DROP COLUMN IF EXISTS direction_type;
