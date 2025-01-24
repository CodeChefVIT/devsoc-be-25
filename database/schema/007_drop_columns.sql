-- +goose Up
ALTER TABLE users DROP COLUMN vit_email;
ALTER TABLE users DROP COLUMN room_no;
ALTER TABLE users DROP COLUMN hostel_block;

-- +goose Down
DROP TABLE users;
