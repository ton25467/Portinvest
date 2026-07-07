-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, role)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password_hash, name, role, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, password_hash, name, role, created_at, updated_at
FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, role, created_at, updated_at
FROM users WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users SET name = $2, role = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
