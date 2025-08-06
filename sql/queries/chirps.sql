-- name: CreateChirp :one
INSERT INTO chirp (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirp :one
SELECT * FROM chirp
WHERE id = $1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirp
WHERE user_id = $1
ORDER BY created_at;

-- name: GetAllChirps :many
SELECT * FROM chirp
ORDER BY created_at;

-- name: DeleteChirp :exec
DELETE FROM chirp
WHERE id = $1 AND user_id = $2;