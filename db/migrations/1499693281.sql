-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE operators DROP COLUMN object_id;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE operators ADD COLUMN object_id text;
