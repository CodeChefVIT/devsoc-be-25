-- +goose Up
ALTER TABLE users ADD COLUMN room_no TEXT DEFAULT NULL;
ALTER TABLE users ADD COLUMN hostel_block TEXT DEFAULT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN room_no;
ALTER TABLE users DROP COLUMN hostel_block;