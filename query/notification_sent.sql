-- name: CreateNotificationSent :one
INSERT INTO notification_sent (
  request_id,
  driver_phone_rejected
) VALUES (
  $1, $2
)
RETURNING *;


-- name: GetNotificationSentByRequestID :one
SELECT * FROM notification_sent
WHERE request_id = $1 LIMIT 1;

-- name: GetNotificationSentByPassengerPhone :one
SELECT * FROM notification_sent JOIN request ON request.id = notification_sent.request_id
WHERE request.phone = $1 LIMIT 1;

-- name: UpdateNotificationSent :one
UPDATE notification_sent
SET driver_phone_rejected = $2
WHERE request_id = $1
RETURNING *;

-- name: DeleteNotificationSent :exec
DELETE FROM notification_sent
WHERE request_id = $1;