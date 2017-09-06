-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas
  ADD COLUMN parent_id uuid,
  ADD COLUMN collect_children boolean,
  ADD COLUMN line_ids text;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas
  DROP COLUMN IF EXISTS parent_id,
  DROP COLUMN IF EXISTS collect_children,
  DROP COLUMN IF EXISTS line_ids;