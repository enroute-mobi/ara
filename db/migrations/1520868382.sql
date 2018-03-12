-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_visits
  DROP COLUMN IF EXISTS arrival_status,
  DROP COLUMN IF EXISTS departure_status,
  DROP COLUMN IF EXISTS collected,
  DROP COLUMN IF EXISTS vehicle_at_stop;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_visits
  ADD COLUMN arrival_status text,
  ADD COLUMN departure_status text,
  ADD COLUMN collected boolean,
  ADD COLUMN vehicle_at_stop boolean;
