-- +goose Up
INSERT INTO ideas (
    id, title, description, track, team_id, is_selected
) VALUES
    ('fc4a4b9d-81ef-4009-bd0a-b678f1a82c93', 'AI Powered Systum', 'An ai system for cleaning machine', 'AI', 'd1a7ffb1-8c85-4d15-b50b-bc120ea97f35', FALSE),
    ('bb3f33a2-bbc9-4b67-9934-9f1d89d44f43', 'Blocdkchain Project', 'A blockchain-based voting system', 'Blockchain', 'a0f4688a-b6fd-4fe2-9d93-0c6de79b98e5', FALSE);

-- +goose Down
DELETE FROM ideas WHERE title IN ('AI Powered App', 'Blockchain Project');
