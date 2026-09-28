-- name: CreateEvent :one
INSERT INTO events (
    provider,
    event_type,
    payload,
    status
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;


-- name: GetEvent :one
SELECT *
FROM events
WHERE id = $1;


-- name: ListEvents :many
SELECT *
FROM events
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;


-- name: UpdateEventStatus :one
UPDATE events
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
