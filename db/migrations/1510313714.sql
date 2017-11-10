-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_areas DROP CONSTRAINT IF EXISTS stop_areas_pkey;
ALTER TABLE stop_areas ADD PRIMARY KEY (id,model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_areas DROP CONSTRAINT IF EXISTS stop_areas_pkey;
ALTER TABLE stop_areas ADD PRIMARY KEY (id);