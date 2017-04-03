-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE lines (
  id              uuid PRIMARY KEY,
  referential_id  uuid NOT NULL,
  name            text,
  object_ids      text,
  attributes      text,
  siri_references text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS lines;