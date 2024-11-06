-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE line_groups (
  id                uuid PRIMARY KEY,
  model_name        text NOT NULL,
  referential_slug  text NOT NULL,
  name              text NOT NULL,
  short_name        text,
  line_ids     text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS line_groups;
