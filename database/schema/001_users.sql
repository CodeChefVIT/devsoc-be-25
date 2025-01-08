-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL UNIQUE,
    team_id UUID,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    phone_no VARCHAR(10) UNIQUE NOT NULL,
    gender CHAR(1) NOT NULL, -- M/F/O
    reg_no TEXT UNIQUE NOT NULL,
    vit_email TEXT UNIQUE NOT NULL,
    hostel_block TEXT NOT NULL,
    room_no INTEGER NOT NULL,
    github_profile TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT not NULL,
    is_leader BOOLEAN NOT NULL,
    is_verified BOOLEAN NOT NULL,
    is_banned BOOLEAN NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
