-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas
  ADD COLUMN referential_id uuid NOT NULL,
  ADD COLUMN requested_at timestamp,
  ADD COLUMN collected_at timestamp,
  ADD COLUMN collected_until timestamp,
  ADD COLUMN collected_always boolean,
  ADD COLUMN object_ids text,
  ADD COLUMN attributes text,
  ADD COLUMN siri_references text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas
  DROP COLUMN IF EXISTS referential_id,
  DROP COLUMN IF EXISTS requested_at,
  DROP COLUMN IF EXISTS collected_at,
  DROP COLUMN IF EXISTS collected_until,
  DROP COLUMN IF EXISTS collected_always,
  DROP COLUMN IF EXISTS object_ids,
  DROP COLUMN IF EXISTS attributes,
  DROP COLUMN IF EXISTS siri_references;