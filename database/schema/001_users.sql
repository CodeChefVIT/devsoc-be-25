-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL UNIQUE,
    name TEXT NOT NULL,
    team_id UUID,
    email TEXT NOT NULL,
    is_vitian BOOLEAN NOT NULL,
    reg_no TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone_no TEXT NOT NULL,
    is_leader BOOLEAN NOT NULL,
    college TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
