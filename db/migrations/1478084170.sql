-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE stop_areas (id UUID, name text);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS stop_areas;