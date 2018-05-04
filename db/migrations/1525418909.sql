-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas ADD COLUMN referent_id uuid;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas DROP COLUMN IF EXISTS referent_id;