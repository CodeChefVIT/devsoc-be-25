-- +goose Up
INSERT INTO users (
    id, name, team_id, email, is_vitian, reg_no, password, phone_no, role, is_leader, college, is_verified, is_banned
) VALUES
    ('9f9bfe51-6b8d-4ff1-8a84-5457a5b8a3c7', 'John Doe', 'd1a7ffb1-8c85-4d15-b50b-bc120ea97f35', 'john.doe@example.com', TRUE, 'VIT2021001', 'password123', '1234567890', 'admin', TRUE, 'VIT', TRUE, FALSE),
    ('0fd44650-776b-11ec-bf63-0242ac130002', 'Jane Doe', 'a0f4688a-b6fd-4fe2-9d93-0c6de79b98e5', 'jane.doe@example.com', TRUE, 'VIT2021002', 'password123', '1234567891', 'student', FALSE, 'VIT', TRUE, FALSE);

-- +goose Down
DELETE FROM users WHERE email IN ('john.doe@example.com', 'jane.doe@example.com');
