-- name: CreateRequest :one
INSERT INTO request (
  type,
  phone,
  pick_up_latitude,
  pick_up_longitude,
  drop_off_latitude,
  drop_off_longitude
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetRequest :one
SELECT * FROM request
WHERE id = $1 LIMIT 1;

-- name: GetRequestByPhone :one
SELECT * FROM request
WHERE phone = $1 LIMIT 1;

-- name: ListRequests :many
SELECT * FROM request
WHERE created_at >= sqlc.arg(start_date)
  AND created_at <= sqlc.arg(end_date)
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateStatusRequest :one
UPDATE request
SET status = $2
WHERE phone = $1
RETURNING *;

-- name: DeleteRequest :exec
DELETE FROM request
WHERE phone = $1;