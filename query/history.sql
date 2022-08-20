-- name: CreateHistory :one
INSERT INTO history (
  type,
  passenger_phone,
  driver_phone,
  pick_up_latitude,
  pick_up_longitude,
  drop_off_latitude,
  drop_off_longitude,
  created_at,
  done_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetHistory :one
SELECT * FROM history
WHERE id = $1 LIMIT 1;

-- name: ListHistoryByPassengerPhone :many
SELECT * FROM history
WHERE passenger_phone = $1;

-- name: ListHistoryByDriverPhone :many
SELECT * FROM history
WHERE driver_phone = $1;

-- name: ListHistorys :many
SELECT * FROM history
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteHistory :exec
DELETE FROM history
WHERE id = $1;