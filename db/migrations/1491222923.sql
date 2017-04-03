-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE ONLY partners ALTER COLUMN slug SET NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE ONLY partners ALTER COLUMN slug DROP NOT NULL;