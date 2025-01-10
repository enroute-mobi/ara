-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE line_groups
  DROP CONSTRAINT line_groups_pkey CASCADE;
ALTER TABLE line_groups
  ADD PRIMARY KEY(id, referential_slug, model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE line_groups
  DROP CONSTRAINT line_groups_pkey CASCADE;
ALTER TABLE line_groups
  ADD PRIMARY KEY(id, model_name);
