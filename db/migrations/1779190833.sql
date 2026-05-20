-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE code_spaces (
    id                uuid PRIMARY KEY,
    referential_id    uuid NOT NULL,
    name              text NOT NULL,
    short_name        text NOT NULL,
    description       text,
    created_at        timestamp NOT NULL,
    updated_at        timestamp NOT NULL
);

CREATE INDEX ON code_spaces (referential_id);
CREATE INDEX ON code_spaces (referential_id, name);
CREATE INDEX ON code_spaces (referential_id, short_name);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS code_spaces;
