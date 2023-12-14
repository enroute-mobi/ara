-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE  lines RENAME COLUMN collect_general_messages TO collect_situations;
ALTER TABLE  stop_areas RENAME COLUMN collect_general_messages TO collect_situations;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE  lines RENAME COLUMN collect_situations TO collect_general_messages;
ALTER TABLE  stop_areas RENAME COLUMN collect_situations TO collect_general_messages;
