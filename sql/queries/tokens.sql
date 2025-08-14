-- name: CreateRefreshToken :one
INSERT INTO tokens ( token, user_id, created_at, updated_at, expires_at, revoked_at )
VALUES ( $1, $2, NOW(), NOW(), $3, NULL )
RETURNING *;

-- name: GetRefreshToken :one
SELECT token, user_id, created_at, updated_at, expires_at, revoked_at
FROM tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
