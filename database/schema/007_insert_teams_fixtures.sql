-- -- +goose Up
-- INSERT INTO teams (
--     id, name, number_of_people, round_qualified, code
-- ) VALUES
--     ('d1a7ffb1-8c85-4d15-b50b-bc120ea97f35', 'Team Alpha', 5, 1, 'ALPHA123'),
--     ('a0f4688a-b6fd-4fe2-9d93-0c6de79b98e5', 'Team Beta', 4, 0, 'BETA456');

-- -- +goose Down
-- DELETE FROM teams WHERE code IN ('ALPHA123', 'BETA456');
