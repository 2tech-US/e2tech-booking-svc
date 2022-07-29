-- name: CreateHistory :one
INSERT INTO history (
  type,
  phone,
  passenger_id,
  driver_id,
  pick_up_latitude,
  pick_up_longitude,
  drop_off_latitude,
  drop_off_longitude,
  created_at,
  done_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetHistory :one
SELECT * FROM history
WHERE id = $1 LIMIT 1;

-- name: ListHistoryByPassengerID :many
SELECT * FROM history
WHERE passenger_id = $1;

-- name: ListHistoryByPhone :many
SELECT * FROM history
WHERE phone = $1;

-- name: ListHistorys :many
SELECT * FROM history
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteHistory :exec
DELETE FROM history
WHERE id = $1;