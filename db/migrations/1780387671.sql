-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE code_spaces
  ADD CONSTRAINT fk_code_spaces_referential_id
  FOREIGN KEY (referential_id) REFERENCES referentials(referential_id)
  DEFERRABLE INITIALLY DEFERRED;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE code_spaces DROP CONSTRAINT fk_code_spaces_referential_id;
