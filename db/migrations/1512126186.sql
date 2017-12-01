-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE lines
  DROP COLUMN IF EXISTS collected_at;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE lines
  ADD COLUMN collected_at timestamp;
