-- name: GetFile :one
SELECT * FROM files
where id = $1 LIMIT 1;

-- name: CreateFile :one
INSERT INTO files (
  id, file_name, extension, content_type,
  file_path, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: CreateAccessLog :one
INSERT INTO access_logs (
  id, ip, method, endpoint, request_body, request_headers,
  request_query, request_params, response_body, response_headers,
  response_time, status_code, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8,
  $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: UpdateAccessLog :one
UPDATE access_logs
SET
  ip = $2,
  method = $3,
  endpoint = $4,
  request_body = $5,
  request_headers = $6,
  request_query = $7,
  request_params = $8,
  response_body = $9,
  response_headers = $10,
  response_time = $11,
  status_code = $12,
  updated_at = $13
WHERE id = $1
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByApiKey :one
SELECT * FROM users
WHERE api_key = $1
LIMIT 1;

-- name: GetUserAndClientByApiKey :one
SELECT
    users.id AS user_id,
    users.username,
    users.api_key,
    users.role AS user_role,
    users.client_id,
    users.created_at AS user_created_at,
    users.updated_at AS user_updated_at,
    users.deleted_at AS user_deleted_at,

    clients.id AS client_id,
    clients.name AS client_name,
    clients.cnpj AS client_cnpj,
    clients.role AS client_role,
    clients.created_at AS client_created_at,
    clients.updated_at AS client_updated_at,
    clients.deleted_at AS client_deleted_at
FROM
    users
JOIN
    clients ON users.client_id = clients.id
WHERE
    users.api_key = $1
LIMIT 1;

-- name: GetAllUsersByClientId :many
SELECT * FROM users
WHERE client_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetClientByUserId :one
SELECT c.* FROM clients c
JOIN users u ON u.client_id = c.id
WHERE u.id = $1
LIMIT 1;

-- name: GetEmissionsByFilters :many
SELECT e.id, e.emission_type, e.client_id, e.message, e.status, e.user_id,
       e.created_at, e.updated_at, e.deleted_at
FROM emissions e
WHERE 
    (sqlc.arg(client_id)::uuid IS NULL OR e.client_id = sqlc.arg(client_id)::uuid)
    AND (sqlc.arg(status)::text IS NULL OR e.status = sqlc.arg(status)::text)
    AND (sqlc.arg(start_date)::timestamp IS NULL OR e.created_at >= sqlc.arg(start_date)::timestamp)
    AND (sqlc.arg(end_date)::timestamp IS NULL OR e.created_at <= sqlc.arg(end_date)::timestamp)
    AND (sqlc.arg(include_deleted)::boolean = TRUE OR e.deleted_at IS NULL)
ORDER BY e.created_at DESC
LIMIT sqlc.arg(row_limit)::integer OFFSET sqlc.arg(row_offset)::integer;

-- name: GetGNREEmissionByChaveNotaAndStatus :one
SELECT g.*
FROM gnre_emissions g
JOIN emissions e ON g.id = e.id
WHERE g.chave_nota = $1 AND e.status = $2
LIMIT 1;

