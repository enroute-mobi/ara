-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE partners DROP CONSTRAINT IF EXISTS partners_slug_key;
ALTER TABLE partners ADD CONSTRAINT slug_referential_unique UNIQUE (slug,referential_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE partners DROP CONSTRAINT IF EXISTS slug_referential_unique;
ALTER TABLE partners ADD CONSTRAINT partners_slug_key UNIQUE (slug);