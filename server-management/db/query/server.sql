-- name: CreateServer :one
INSERT INTO servers (
  name, 
  ipv4, 
  status
) VALUES (
  @name::text, @ipv4::text, @status::int
)
RETURNING *;

-- name: GetServer :many
SELECT * FROM servers
WHERE true
    AND (CASE WHEN @is_id::bool THEN id = @id::int ELSE TRUE END)
    AND (CASE WHEN @is_name::bool THEN name = @name::text ELSE TRUE END) -- :if @name
    AND (CASE WHEN @is_status::bool THEN status = @status::int ELSE TRUE END) -- :if @status
    AND (CASE WHEN @is_ipv4::bool THEN ipv4 = @ipv4::text ELSE TRUE END) -- :if @ipv4
ORDER BY
    CASE WHEN @id_asc::bool THEN id END ASC, -- :if @id_asc
    CASE WHEN @id_desc::bool  THEN id END DESC, -- :if @id_desc
    CASE WHEN @name_asc::bool  THEN name END ASC, -- :if @name_asc
    CASE WHEN @name_desc::bool  THEN name END DESC, -- :if @name_desc
    CASE WHEN @status_asc::bool  THEN status END ASC, -- :if @status_asc
    CASE WHEN @status_desc::bool  THEN status END DESC, -- :if @status_desc
    CASE WHEN @ipv4_asc::bool  THEN ipv4 END ASC, -- :if @ipv4_asc
    CASE WHEN @ipv4_desc::bool  THEN ipv4 END DESC, -- :if @ipv4_desc
    CASE WHEN @created_at_asc::bool  THEN created_at END ASC, -- :if @created_at_asc
    CASE WHEN @created_at_desc::bool  THEN created_at END DESC, -- :if @created_at_desc
    CASE WHEN @updated_at_asc::bool  THEN update_at END ASC, -- :if @updated_at_asc
    CASE WHEN @updated_at_desc::bool  THEN update_at END DESC -- :if @updated_at_desc
LIMIT $1
OFFSET $2;

-- name: UpdateServer :one
UPDATE servers 
SET 
    name = CASE WHEN @set_name::bool THEN @name::text ELSE name END, -- :if @name
    status = CASE WHEN @set_status::bool THEN @status::int ELSE status END, -- :if @status
    ipv4 = CASE WHEN @set_ipv4::bool THEN @ipv4::text ELSE ipv4 END, -- :if @ipv4
    update_at = CASE WHEN @set_update_at::bool THEN @update_at::timestamp ELSE update_at END -- :if @update_at
WHERE id = $1
RETURNING *; 

-- name: DeleteServer :exec
DELETE FROM servers 
WHERE id = $1;

-- name: GetAllServers :many
SELECT * FROM servers;