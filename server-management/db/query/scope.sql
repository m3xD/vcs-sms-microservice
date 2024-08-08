-- name: CreateScope :one
INSERT INTO scopes ("name", "role")
VALUES ($1, $2)
RETURNING *;

-- name: GetScope :many
SELECT S.name FROM scopes S
JOIN users U ON U.role = S.role
WHERE U.id = $1;

-- name: UpdateScope :one
UPDATE scopes
SET name = name = CASE WHEN @set_name::bool THEN @name::text ELSE name END,
    role = CASE WHEN @set_role::bool THEN @role::text ELSE role END
WHERE id = $1
RETURNING *;

-- name: DeleteScope :exec 
DELETE FROM scopes
WHERE id = $1;