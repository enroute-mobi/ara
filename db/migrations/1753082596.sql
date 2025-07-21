-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE facilities (
id                uuid NOT NULL,
referential_slug  text NOT NULL,
model_name        text NOT NULL,
codes             text,
PRIMARY KEY (id, referential_slug, model_name)
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS facilities;
