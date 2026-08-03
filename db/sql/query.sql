-- name: CreateUser :exec
WITH new_user AS (
    INSERT INTO users (username, email, is_active)
    VALUES ($1, $2, true)
    RETURNING id
)
INSERT INTO user_pass (id, password_hash)
SELECT id, $3
FROM new_user;

-- name: GetUserFromName :one
SELECT id, username, email, is_active, created_at
FROM users
WHERE username = $1;

-- name: GetUserFromID :one
SELECT id, username, email, is_active, created_at
FROM users
WHERE id = $1;

-- name: GetHash :one
SELECT password_hash
FROM user_pass
WHERE id = $1;

-- name: UpdateHash :exec
UPDATE user_pass
SET password_hash = $2
WHERE id = $1;
