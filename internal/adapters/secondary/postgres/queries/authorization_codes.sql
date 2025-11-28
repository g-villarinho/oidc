-- name: CreateAuthorizationCode :one
INSERT INTO authorization_codes (
    code,
    client_id,
    user_id,
    redirect_uri,
    scopes,
    code_challenge,
    code_challenge_method,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAuthorizationCode :one
SELECT 
    ac.*,
    c.client_id as client_client_id,
    c.redirect_uris as client_redirect_uris,
    u.email as user_email
FROM authorization_codes ac
JOIN oauth_clients c ON ac.client_id = c.client_id
JOIN users u ON ac.user_id = u.id
WHERE ac.code = $1 
  AND ac.expires_at > NOW()
LIMIT 1;

-- name: DeleteAuthorizationCode :exec
DELETE FROM authorization_codes
WHERE code = $1;

-- name: DeleteExpiredAuthorizationCodes :exec
DELETE FROM authorization_codes
WHERE expires_at < NOW();
