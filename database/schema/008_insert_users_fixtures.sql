-- +goose Up
INSERT INTO users (
    id, 
    team_id, 
    first_name,
    last_name,
    email,
    phone_no,
    gender,
    reg_no,
    vit_email,
    hostel_block,
    room_no,
    github_profile,
    password,
    role,
    is_leader,
    is_verified,
    is_banned
) VALUES
    (
        '9f9bfe51-6b8d-4ff1-8a84-5457a5b8a3c7',
        'd1a7ffb1-8c85-4d15-b50b-bc120ea97f35',
        'John',
        'Doe',
        'john.doe@example.com',
        '1234567890',
        'M',
        'VIT2021001',
        'john.doe@vitstudent.ac.in',
        'A',
        101,
        'https://github.com/johndoe',
        'password123',
        'admin',
        TRUE,
        TRUE,
        FALSE
    ),
    (
        '0fd44650-776b-11ec-bf63-0242ac130002',
        'a0f4688a-b6fd-4fe2-9d93-0c6de79b98e5',
        'Jane',
        'Doe', 
        'jane.doe@example.com',
        '1234567891',
        'F',
        'VIT2021002',
        'jane.doe@vitstudent.ac.in',
        'B',
        202,
        'https://github.com/janedoe',
        'password123',
        'student',
        FALSE,
        TRUE,
        FALSE
    );

-- +goose Down
DELETE FROM users WHERE email IN ('john.doe@example.com', 'jane.doe@example.com');
