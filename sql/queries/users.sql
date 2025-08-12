-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_passwords)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM users 
WHERE id IN (
    SELECT user_id
    FROM refresh_tokens
    WHERE token = $1
    AND revoked_at IS NULL
    AND expires_at > NOW()
);

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_passwords = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;