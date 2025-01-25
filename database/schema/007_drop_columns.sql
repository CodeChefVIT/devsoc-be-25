-- +goose Up
ALTER TABLE users DROP COLUMN vit_email;
ALTER TABLE users DROP COLUMN room_no;
ALTER TABLE users DROP COLUMN hostel_block;

-- +goose Down
ALTER TABLE users ADD COLUMN vit_email TEXT;
ALTER TABLE users ADD COLUMN room_no INTEGER DEFAULT NULL;
ALTER TABLE users ADD COLUMN hostel_block TEXT DEFAULT NULL;