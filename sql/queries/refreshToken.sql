-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING *;


-- name: GetUserFromRefreshToken :one
SELECT u.id, u.email, u.created_at, u.updated_at, u.hashed_password, rt.expires_at, rt.revoked_at
FROM refresh_tokens rt
JOIN users u ON rt.user_id = u.id
WHERE rt.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

