-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE partners ADD COLUMN name text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE partners DROP COLUMN IF EXISTS name;