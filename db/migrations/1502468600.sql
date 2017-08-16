-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas ADD COLUMN model_name text NOT NULL;
ALTER TABLE lines ADD COLUMN model_name text NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas DROP COLUMN IF EXISTS model;
ALTER TABLE lines DROP COLUMN IF EXISTS model;
