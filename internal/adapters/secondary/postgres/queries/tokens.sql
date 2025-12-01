-- name: CreateToken :one
INSERT INTO tokens (
    id,
    access_token_hash,
    refresh_token_hash,
    authorization_code,
    client_id,
    user_id,
    scopes,
    token_type,
    access_token_expires_at,
    refresh_token_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetTokenByAccessTokenHash :one
SELECT * FROM tokens
WHERE access_token_hash = $1
  AND revoked = FALSE
LIMIT 1;

-- name: GetTokenByRefreshTokenHash :one
SELECT * FROM tokens
WHERE refresh_token_hash = $1
  AND revoked = FALSE
  AND refresh_token_expires_at > NOW()
LIMIT 1;

-- name: GetTokenByID :one
SELECT * FROM tokens
WHERE id = $1
LIMIT 1;

-- name: RevokeToken :exec
UPDATE tokens
SET
    revoked = TRUE,
    revoked_at = NOW(),
    revoked_reason = $2
WHERE id = $1;

-- name: RevokeTokenByAccessTokenHash :exec
UPDATE tokens
SET
    revoked = TRUE,
    revoked_at = NOW(),
    revoked_reason = $2
WHERE access_token_hash = $1;

-- name: RevokeTokensByUser :exec
UPDATE tokens
SET
    revoked = TRUE,
    revoked_at = NOW(),
    revoked_reason = $2
WHERE user_id = $1
  AND revoked = FALSE;

-- name: RevokeTokensByClient :exec
UPDATE tokens
SET
    revoked = TRUE,
    revoked_at = NOW(),
    revoked_reason = $2
WHERE client_id = $1
  AND revoked = FALSE;

-- name: RevokeTokensByAuthorizationCode :exec
UPDATE tokens
SET
    revoked = TRUE,
    revoked_at = NOW(),
    revoked_reason = $2
WHERE authorization_code = $1
  AND revoked = FALSE;

-- name: UpdateLastUsedAt :exec
UPDATE tokens
SET last_used_at = NOW()
WHERE id = $1;

-- name: UpdateLastUsedAtByAccessTokenHash :exec
UPDATE tokens
SET last_used_at = NOW()
WHERE access_token_hash = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM tokens
WHERE access_token_expires_at < NOW()
  AND refresh_token_expires_at < NOW();

-- name: GetActiveTokensByUser :many
SELECT * FROM tokens
WHERE user_id = $1
  AND revoked = FALSE
  AND access_token_expires_at > NOW()
ORDER BY created_at DESC;

-- name: GetActiveTokensByClient :many
SELECT * FROM tokens
WHERE client_id = $1
  AND revoked = FALSE
  AND access_token_expires_at > NOW()
ORDER BY created_at DESC;

-- name: GetTokenWithDetails :one
SELECT
    t.*,
    u.email as user_email,
    u.name as user_name,
    c.client_name as client_name
FROM tokens t
JOIN users u ON t.user_id = u.id
JOIN oauth_clients c ON t.client_id = c.client_id
WHERE t.id = $1
LIMIT 1;

-- name: CountActiveTokensByUser :one
SELECT COUNT(*) FROM tokens
WHERE user_id = $1
  AND revoked = FALSE
  AND access_token_expires_at > NOW();

-- name: CountActiveTokensByClient :one
SELECT COUNT(*) FROM tokens
WHERE client_id = $1
  AND revoked = FALSE
  AND access_token_expires_at > NOW();
