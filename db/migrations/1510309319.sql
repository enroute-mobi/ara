-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_visits DROP COLUMN IF EXISTS referential_id;
ALTER TABLE stop_visits ADD COLUMN referential_slug text NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_visits DROP COLUMN IF EXISTS referential_slug;
ALTER TABLE stop_visits ADD COLUMN referential_id uuid NOT NULL;