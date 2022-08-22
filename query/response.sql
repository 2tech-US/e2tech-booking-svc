-- name: CreateResponse :one
INSERT INTO response (
  request_id,
  driver_phone,
  driver_name,
  driver_latitude,
  driver_longitude
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetResponse :one
SELECT * FROM response
WHERE id = $1 LIMIT 1;

-- name: GetResponseByRequestID :one
SELECT * FROM response
WHERE request_id = $1 LIMIT 1;

-- name: GetResponseByPassengerPhone :one
SELECT * FROM response JOIN request ON request.id = response.request_id
WHERE request.phone = $1 LIMIT 1;

-- name: GetResponseByDriverPhone :one
SELECT * FROM response 
WHERE driver_phone = $1 LIMIT 1;

-- name: ListResponses :many
SELECT * FROM response
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateResponse :one
UPDATE response
SET driver_phone = $2,
    driver_latitude = $3,
    driver_longitude = $4
WHERE id = $1
RETURNING *;

-- name: DeleteResponse :exec
DELETE FROM response
WHERE driver_phone = $1;

-- name: DeleteResponseByRequestID :exec
DELETE FROM response
WHERE request_id = $1;