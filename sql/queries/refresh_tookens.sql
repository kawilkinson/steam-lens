-- #nosec G101 -- token is dynamic and not a hardcoded credential
-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- #nosec G101 -- false positive, token is a SQL placeholder
-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE token = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- #nosec G101 -- false positive, token is just used in SQL update
-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
