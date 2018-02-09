-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE operators DROP CONSTRAINT IF EXISTS operators_pkey;
ALTER TABLE operators ADD PRIMARY KEY (id,model_name);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE operators DROP CONSTRAINT IF EXISTS operators_pkey;
ALTER TABLE operators ADD PRIMARY KEY (id);