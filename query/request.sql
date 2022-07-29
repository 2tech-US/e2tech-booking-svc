-- name: CreateRequest :one
INSERT INTO request (
  type,
  passenger_id,
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

-- name: GetRequestByPassengerID :one
SELECT * FROM request
WHERE passenger_id = $1 LIMIT 1;

-- name: ListRequests :many
SELECT * FROM request
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateStatusRequest :one
UPDATE request
SET status = $2
WHERE id = $1
RETURNING *;

-- name: DeleteRequest :exec
DELETE FROM request
WHERE id = $1;