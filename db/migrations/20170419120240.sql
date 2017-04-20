-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE  ONLY  stop_areas RENAME COLUMN requested_at TO next_collect_at;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE  ONLY  stop_areas RENAME COLUMN next_collect_at TO requested_at;
