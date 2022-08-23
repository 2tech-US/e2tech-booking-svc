-- name: CreateHistory :one
INSERT INTO history (
    type,
    passenger_phone,
    driver_phone,
    pick_up_latitude,
    pick_up_longitude,
    drop_off_latitude,
    drop_off_longitude,
    price,
    created_at,
    done_at
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
-- name: GetHistory :one
SELECT *
FROM history
WHERE id = $1
LIMIT 1;
-- name: ListHistoryByPassengerPhone :many
SELECT *
FROM history
WHERE passenger_phone = $1
  AND created_at >= sqlc.arg(start_date)
  AND created_at <= sqlc.arg(end_date)
ORDER BY id
LIMIT $2 OFFSET $3;
-- name: ListHistoryByDriverPhone :many
SELECT *
FROM history
WHERE driver_phone = $1
  AND created_at >= sqlc.arg(start_date)
  AND created_at <= sqlc.arg(end_date)
ORDER BY id
LIMIT $2 OFFSET $3;
-- name: ListHistories :many
SELECT *
FROM history
WHERE created_at >= sqlc.arg(start_date)
  AND created_at <= sqlc.arg(end_date)
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: DeleteHistory :exec
DELETE FROM history
WHERE id = $1;