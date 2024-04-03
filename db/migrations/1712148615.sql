-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE lines ADD COLUMN referent_id uuid;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE lines DROP COLUMN IF EXISTS referent_id;