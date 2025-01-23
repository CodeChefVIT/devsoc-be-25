-- +goose up
ALTER TABLE users
DROP CONSTRAINT users_vit_email_key;


ALTER TABLE users
DROP CONSTRAINT users_phone_no_key;

ALTER TABLE users
DROP CONSTRAINT users_reg_no_key;

-- +goose Down
-- This migration is irreversible and does not perform any action on rollback.
