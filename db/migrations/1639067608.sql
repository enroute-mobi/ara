-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE operators DROP CONSTRAINT operators_pkey;
ALTER TABLE operators ADD PRIMARY KEY (id, referential_slug, model_name);
CREATE INDEX ON operators (referential_slug, model_name);

ALTER TABLE stop_areas DROP CONSTRAINT stop_areas_pkey;
ALTER TABLE stop_areas ADD PRIMARY KEY (id, referential_slug, model_name);
CREATE INDEX ON stop_areas (referential_slug, model_name);

ALTER TABLE lines DROP CONSTRAINT lines_pkey;
ALTER TABLE lines ADD PRIMARY KEY (id, referential_slug, model_name);
CREATE INDEX ON lines (referential_slug, model_name);

ALTER TABLE vehicle_journeys DROP CONSTRAINT vehicle_journeys_pkey;
ALTER TABLE vehicle_journeys ADD PRIMARY KEY (id, referential_slug, model_name);
CREATE INDEX ON vehicle_journeys (referential_slug, model_name);

ALTER TABLE stop_visits DROP CONSTRAINT stop_visits_pkey;
ALTER TABLE stop_visits ADD PRIMARY KEY (id, referential_slug, model_name);
CREATE INDEX ON stop_visits (referential_slug, model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE operators DROP CONSTRAINT IF EXISTS operators_pkey;
ALTER TABLE operators ADD PRIMARY KEY (id, model_name);
DROP INDEX IF EXISTS operators_referential_slug_model_name_idx;

ALTER TABLE stop_areas DROP CONSTRAINT IF EXISTS stop_areas_pkey;
ALTER TABLE stop_areas ADD PRIMARY KEY (id, model_name);
DROP INDEX IF EXISTS stop_areas_referential_slug_model_name_idx;

ALTER TABLE lines DROP CONSTRAINT IF EXISTS lines_pkey;
ALTER TABLE lines ADD PRIMARY KEY (id, model_name);
DROP INDEX IF EXISTS lines_referential_slug_model_name_idx;

ALTER TABLE vehicle_journeys DROP CONSTRAINT IF EXISTS vehicle_journeys_pkey;
ALTER TABLE vehicle_journeys ADD PRIMARY KEY (id, model_name);
DROP INDEX IF EXISTS vehicle_journeys_referential_slug_model_name_idx;

ALTER TABLE stop_visits DROP CONSTRAINT IF EXISTS stop_visits_pkey;
ALTER TABLE stop_visits ADD PRIMARY KEY (id, model_name);
DROP INDEX IF EXISTS stop_visits_referential_slug_model_name_idx;
