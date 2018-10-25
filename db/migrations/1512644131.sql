-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE operators ADD model_name text NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE operators DROP model_name;

