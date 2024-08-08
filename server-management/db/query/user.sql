-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email,
  role
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: UpdateRole :one
UPDATE users
SET role = $1
WHERE id = $2
RETURNING *;