-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    name,
    email_verified
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;


-- name: GetByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    email = $2,
    password_hash = $3,
    name = $4,
    email_verified = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;    

-- name: VerifyEmail :one
UPDATE users
SET 
    email_verified = TRUE,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePassword :one
UPDATE users
SET 
    password_hash = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;