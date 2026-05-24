-- name: CreateWes :one
INSERT INTO wes (height)
VALUES ($1)
RETURNING *;

-- name: CreateMoves :many
INSERT INTO moves (wes_id, status, description)
SELECT $1::uuid, UNNEST($2::move_status[]), UNNEST($3::text[])
RETURNING *;

-- name: GetWesByID :one
SELECT *
FROM wes
WHERE id = $1;

-- name: GetAllMovesByWesID :many
SELECT *
FROM moves
WHERE wes_id = $1
ORDER BY created_at;
