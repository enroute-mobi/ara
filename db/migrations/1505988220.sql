-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE lines
  ADD COLUMN collect_general_messages boolean,
  ADD COLUMN collected_at timestamp;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE lines
  DROP COLUMN IF EXISTS collect_general_messages,
  DROP COLUMN IF EXISTS collected_at;