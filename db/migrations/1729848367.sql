-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE stop_area_groups (
  id                uuid PRIMARY KEY,
  model_name        text NOT NULL,
  referential_slug  text NOT NULL,
  name              text NOT NULL,
  short_name        text,
  stop_area_ids     text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS stop_area_groups;
