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

-- name: GetGNREEmissionById :one
SELECT g.*
FROM gnre_emissions g
JOIN emissions e ON g.id = e.id
WHERE g.id = $1
LIMIT 1;

-- name: CreateGNREEmission :one
WITH new_emission AS (
    INSERT INTO emissions (
                           id,
                           emission_type,
                           client_id,
                           message,
                           status,
                           user_id
        ) VALUES (
                     gen_random_uuid(),
                     sqlc.arg(emission_type)::text,
                     sqlc.arg(client_id)::uuid,
                     COALESCE(sqlc.arg(message)::text, ''),
                     COALESCE(sqlc.arg(status)::text, 'PENDING'),
                     sqlc.arg(user_id)::uuid
                 ) RETURNING id, created_at
)
INSERT INTO gnre_emissions (
    id,
    xml,
    guia_amount,
    chave_nota,
    num_nota,
    destinatario
) SELECT
      ne.id,
      sqlc.arg(xml_file_id)::uuid,
      sqlc.arg(guia_amount)::double precision,
      sqlc.arg(chave_nota)::varchar,
      sqlc.arg(num_nota)::varchar,
      sqlc.arg(destinatario)::varchar
FROM new_emission ne
RETURNING
    id,
    (SELECT emission_type FROM emissions WHERE id = gnre_emissions.id) AS emission_type,
    (SELECT client_id FROM emissions WHERE id = gnre_emissions.id) AS client_id,
    (SELECT COALESCE(message, '') FROM emissions WHERE id = gnre_emissions.id) AS message,
    (SELECT COALESCE(status, 'PENDING') FROM emissions WHERE id = gnre_emissions.id) AS status,
    (SELECT user_id FROM emissions WHERE id = gnre_emissions.id) AS user_id,
    (SELECT created_at FROM emissions WHERE id = gnre_emissions.id) AS created_at,
    xml AS xml_file_id,
    guia_amount,
    chave_nota,
    num_nota,
    destinatario;

-- name: UpdateEmissionAndGNRE :one
WITH check_exists AS (
    SELECT e.id
    FROM emissions e
             JOIN gnre_emissions g ON e.id = g.id
    WHERE e.id = sqlc.arg(id)
    LIMIT 1
),
     update_emission AS (
         UPDATE emissions
             SET
                 emission_type = COALESCE(sqlc.narg('emission_type'), emission_type),
                 client_id = COALESCE(sqlc.narg('client_id'), client_id),
                 message = COALESCE(sqlc.narg('message'), message),
                 status = COALESCE(sqlc.narg('status'), status),
                 user_id = COALESCE(sqlc.narg('user_id'), user_id),
                 updated_at = CURRENT_TIMESTAMP
             WHERE id = (SELECT id FROM check_exists)
             RETURNING *
     ),
     update_gnre AS (
         UPDATE gnre_emissions
             SET
                 xml = COALESCE(sqlc.narg('xml')::uuid, xml),
                 pdf = COALESCE(sqlc.narg('pdf')::uuid, pdf),
                 comprovante_pdf = COALESCE(sqlc.narg('comprovante_pdf')::uuid, comprovante_pdf),
                 guia_amount = COALESCE(sqlc.narg('guia_amount')::numeric, guia_amount),
                 numero_recibo = COALESCE(sqlc.narg('numero_recibo')::varchar, numero_recibo),
                 chave_nota = COALESCE(sqlc.narg('chave_nota')::varchar, chave_nota),
                 cod_barras_guia = COALESCE(sqlc.narg('cod_barras_guia')::varchar, cod_barras_guia),
                 num_nota = COALESCE(sqlc.narg('num_nota')::varchar, num_nota),
                 destinatario = COALESCE(sqlc.narg('destinatario')::varchar, destinatario)
             WHERE id = (SELECT id FROM check_exists)
             RETURNING *
     )
SELECT
    e.id AS emission_id,
    e.emission_type,
    e.client_id,
    e.message,
    e.status,
    e.user_id,
    e.created_at AS emission_created_at,
    e.updated_at AS emission_updated_at,
    e.deleted_at AS emission_deleted_at,
    g.xml AS gnre_xml,
    g.pdf AS gnre_pdf,
    g.comprovante_pdf AS gnre_comprovante_pdf,
    g.guia_amount AS gnre_guia_amount,
    g.numero_recibo AS gnre_numero_recibo,
    g.chave_nota AS gnre_chave_nota,
    g.cod_barras_guia AS gnre_cod_barras_guia,
    g.num_nota AS gnre_num_nota,
    g.destinatario AS gnre_destinatario
FROM update_emission e
         JOIN update_gnre g ON e.id = g.id;
