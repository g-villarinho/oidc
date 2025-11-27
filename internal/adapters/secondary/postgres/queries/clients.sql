-- name: GetClientByClientID :one
SELECT * FROM oauth_clients
WHERE client_id = $1 LIMIT 1;

-- name: GetClientByID :one
SELECT * FROM oauth_clients
WHERE id = $1 LIMIT 1;

-- name: CreateClient :one
INSERT INTO oauth_clients (
    id,
    client_id,
    client_secret,
    client_name,
    redirect_uris,
    grant_types,
    response_types,
    scopes,
    logo_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: ListClients :many
SELECT * FROM oauth_clients
ORDER BY created_at DESC;

-- name: UpdateClient :one
UPDATE oauth_clients
SET 
    client_name = $2,
    redirect_uris = $3,
    grant_types = $4,
    response_types = $5,
    scopes = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM oauth_clients
WHERE id = $1;