-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE facilities RENAME COLUMN model_name TO model_date;
ALTER TABLE line_groups RENAME COLUMN model_name TO model_date;
ALTER TABLE lines RENAME COLUMN model_name TO model_date;
ALTER TABLE operators RENAME COLUMN model_name TO model_date;
ALTER TABLE stop_area_groups RENAME COLUMN model_name TO model_date;
ALTER TABLE stop_areas RENAME COLUMN model_name TO model_date;
ALTER TABLE stop_visits RENAME COLUMN model_name TO model_date;
ALTER TABLE vehicle_journeys RENAME COLUMN model_name TO model_date;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE facilities RENAME COLUMN model_date TO model_name;
ALTER TABLE line_groups RENAME COLUMN model_date TO model_name;
ALTER TABLE lines RENAME COLUMN model_date TO model_name;
ALTER TABLE operators RENAME COLUMN model_date TO model_name;
ALTER TABLE stop_area_groups RENAME COLUMN model_date TO model_name;
ALTER TABLE stop_areas RENAME COLUMN model_date TO model_name;
ALTER TABLE stop_visits RENAME COLUMN model_date TO model_name;
ALTER TABLE vehicle_journeys RENAME COLUMN model_date TO model_name;
