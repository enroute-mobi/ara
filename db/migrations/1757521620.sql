-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE partner_templates (
    id                uuid NOT NULL,
    referential_id    uuid NOT NULL,
    slug              text UNIQUE,
    credential_type   text NOT NULL,
    local_credential  text NOT NULL,
    remote_credential text NOT NULL,
    max_partners      smallint,
    name              text,
    settings          text,
    connector_types   text
);

ALTER TABLE ONLY partner_templates ADD CONSTRAINT partner_templates_pkey PRIMARY KEY (id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS partner_templates;
