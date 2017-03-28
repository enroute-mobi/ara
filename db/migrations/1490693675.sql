-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE ONLY stop_areas ADD CONSTRAINT stop_areas_pkey PRIMARY KEY (id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE ONLY stop_areas DROP CONSTRAINT IF EXISTS stop_areas_pkey;