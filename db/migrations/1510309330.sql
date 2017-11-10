-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE vehicle_journeys DROP COLUMN IF EXISTS referential_id;
ALTER TABLE vehicle_journeys ADD COLUMN referential_slug text NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE vehicle_journeys DROP COLUMN IF EXISTS referential_slug;
ALTER TABLE vehicle_journeys ADD COLUMN referential_id uuid NOT NULL;