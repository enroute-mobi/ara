-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE macros (
  id               uuid PRIMARY KEY,
  referential_slug text NOT NULL,
  context_id       text,
  position         smallint,
  type             text NOT NULL,
  model_type       text,
  hook             text,
  attributes       text
);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS macros;