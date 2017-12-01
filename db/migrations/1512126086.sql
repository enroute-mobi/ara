-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas
  DROP COLUMN IF EXISTS next_collect_at,
  DROP COLUMN IF EXISTS collected_at,
  DROP COLUMN IF EXISTS collected_until;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas
  ADD COLUMN next_collect_at timestamp,
  ADD COLUMN collected_at timestamp,
  ADD COLUMN collected_until timestamp;
