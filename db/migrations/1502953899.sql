-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE stop_visits (
  id                  uuid PRIMARY KEY,
  model_name          text NOT NULL,
  referential_id      uuid NOT NULL,
  object_ids          text,
  stop_area_id        uuid,
  vehicle_journey_id  uuid,
  collected           boolean,
  vehicle_at_stop     boolean,
  collected_at        timestamp,
  recorded_at         timestamp,
  passage_order       smallint,
  arrival_status      text,
  departure_status    text,
  schedules           text,
  attributes          text,
  siri_references     text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS stop_visits;
