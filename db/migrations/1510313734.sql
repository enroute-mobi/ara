-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE vehicle_journeys DROP CONSTRAINT IF EXISTS vehicle_journeys_pkey;
ALTER TABLE vehicle_journeys ADD PRIMARY KEY (id,model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE vehicle_journeys DROP CONSTRAINT IF EXISTS vehicle_journeys_pkey;
ALTER TABLE vehicle_journeys ADD PRIMARY KEY (id);