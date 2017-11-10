-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE stop_visits DROP CONSTRAINT IF EXISTS stop_visits_pkey;
ALTER TABLE stop_visits ADD PRIMARY KEY (id,model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE stop_visits DROP CONSTRAINT IF EXISTS stop_visits_pkey;
ALTER TABLE stop_visits ADD PRIMARY KEY (id);