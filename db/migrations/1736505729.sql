-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE stop_area_groups
  DROP CONSTRAINT stop_area_groups_pkey CASCADE;
ALTER TABLE stop_area_groups
  ADD PRIMARY KEY (id, referential_slug, model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE stop_area_groups
  DROP CONSTRAINT stop_area_groups_pkey CASCADE;
ALTER TABLE stop_area_groups
  ADD PRIMARY KEY (id, model_name);
