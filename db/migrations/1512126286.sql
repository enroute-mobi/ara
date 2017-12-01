-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_visits
  DROP COLUMN IF EXISTS collected_at,
  DROP COLUMN IF EXISTS recorded_at;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_visits
  ADD COLUMN collected_at timestamp,
  ADD COLUMN recorded_at timestamp;
