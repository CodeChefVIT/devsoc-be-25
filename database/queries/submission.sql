-- name: GetSubmissionByTeamID :one
SELECT * FROM submission WHERE team_id = $1;

-- name: CreateSubmission :one
INSERT INTO submission (
    id,
    team_id,
    title,
    description,
    github_link,
    figma_link,
    ppt_link,
    other_link
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateSubmission :one
UPDATE submission
SET github_link = $2,
    figma_link = $3,
    ppt_link = $4,
    other_link = $5,
    title = $6,
    description = $7
WHERE team_id = $1
RETURNING *;

-- name: DeleteSubmission :exec
DELETE FROM submission WHERE team_id = $1;
