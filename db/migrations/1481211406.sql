-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE partners (
    id              uuid NOT NULL,
    referential_id  uuid NOT NULL,
    slug            text UNIQUE,
    settings        text,
    connector_types text
);

ALTER TABLE ONLY partners ADD CONSTRAINT partners_pkey PRIMARY KEY (id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS partners;
