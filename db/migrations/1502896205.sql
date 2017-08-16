-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE vehicle_journeys (
  id              uuid PRIMARY KEY,
  model_name      text NOT NULL,
  referential_id  uuid NOT NULL,
  name            text,
  object_ids      text,
  line_id         uuid,
  attributes      text,
  siri_references text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS vehicle_journeys;
