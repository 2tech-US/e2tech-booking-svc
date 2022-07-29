-- name: CreateResponse :one
INSERT INTO response (
  request_id,
  driver_id,
  driver_latitude,
  driver_longitude
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetResponse :one
SELECT * FROM response
WHERE id = $1 LIMIT 1;

-- name: GetResponseByRequestID :one
SELECT * FROM response
WHERE request_id = $1 LIMIT 1;

-- name: GetResponseByPassengerID :one
SELECT * FROM response
WHERE request_id = $1 LIMIT 1;

-- name: ListResponses :many
SELECT * FROM response
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateResponse :one
UPDATE response
SET driver_id = $2,
    driver_latitude = $3,
    driver_longitude = $4
WHERE id = $1
RETURNING *;

-- name: DeleteResponse :exec
DELETE FROM response
WHERE id = $1;