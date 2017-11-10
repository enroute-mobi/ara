-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE lines DROP CONSTRAINT IF EXISTS lines_pkey;
ALTER TABLE lines ADD PRIMARY KEY (id,model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE lines DROP CONSTRAINT IF EXISTS lines_pkey;
ALTER TABLE lines ADD PRIMARY KEY (id);