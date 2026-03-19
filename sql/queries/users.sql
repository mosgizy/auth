-- name: CreateUser :one
INSERT INTO users (email,password_hash,full_name)
VALUES ($1,$2,$3)
RETURNING id, email, full_name, created_at;

-- name: GetUserByEmail :one
SELECT id, email, full_name, password_hash
FROM users
WHERE email = $1;

-- name: SetResetToken :exec
UPDATE users
SET reset_token = $2, reset_token_expires = $3
WHERE email = $1;

-- name: GetUserByResetToken :one
SELECT * FROM users
WHERE reset_token = $1 AND reset_token_expires > NOW();

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = $2, 
  reset_token = NULL, 
  reset_token_expires = NULL, 
  updated_at = NOW()
WHERE id = $1;