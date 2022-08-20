// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: notification_sent.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createNotificationSent = `-- name: CreateNotificationSent :one
INSERT INTO notification_sent (
  request_id,
  driver_phone_rejected
) VALUES (
  $1, $2
)
RETURNING id, request_id, driver_phone_rejected
`

type CreateNotificationSentParams struct {
	RequestID           int64    `json:"request_id"`
	DriverPhoneRejected []string `json:"driver_phone_rejected"`
}

func (q *Queries) CreateNotificationSent(ctx context.Context, arg CreateNotificationSentParams) (NotificationSent, error) {
	row := q.db.QueryRowContext(ctx, createNotificationSent, arg.RequestID, pq.Array(arg.DriverPhoneRejected))
	var i NotificationSent
	err := row.Scan(&i.ID, &i.RequestID, pq.Array(&i.DriverPhoneRejected))
	return i, err
}

const deleteNotificationSent = `-- name: DeleteNotificationSent :exec
DELETE FROM notification_sent
WHERE request_id = $1
`

func (q *Queries) DeleteNotificationSent(ctx context.Context, requestID int64) error {
	_, err := q.db.ExecContext(ctx, deleteNotificationSent, requestID)
	return err
}

const getNotificationSentByPassengerPhone = `-- name: GetNotificationSentByPassengerPhone :one
SELECT notification_sent.id, request_id, driver_phone_rejected, request.id, type, phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, status, created_at, expire_at FROM notification_sent JOIN request ON request.id = notification_sent.request_id
WHERE request.phone = $1 LIMIT 1
`

type GetNotificationSentByPassengerPhoneRow struct {
	ID                  int64        `json:"id"`
	RequestID           int64        `json:"request_id"`
	DriverPhoneRejected []string     `json:"driver_phone_rejected"`
	ID_2                int64        `json:"id_2"`
	Type                string       `json:"type"`
	Phone               string       `json:"phone"`
	PickUpLatitude      float64      `json:"pick_up_latitude"`
	PickUpLongitude     float64      `json:"pick_up_longitude"`
	DropOffLatitude     float64      `json:"drop_off_latitude"`
	DropOffLongitude    float64      `json:"drop_off_longitude"`
	Status              string       `json:"status"`
	CreatedAt           time.Time    `json:"created_at"`
	ExpireAt            sql.NullTime `json:"expire_at"`
}

func (q *Queries) GetNotificationSentByPassengerPhone(ctx context.Context, phone string) (GetNotificationSentByPassengerPhoneRow, error) {
	row := q.db.QueryRowContext(ctx, getNotificationSentByPassengerPhone, phone)
	var i GetNotificationSentByPassengerPhoneRow
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		pq.Array(&i.DriverPhoneRejected),
		&i.ID_2,
		&i.Type,
		&i.Phone,
		&i.PickUpLatitude,
		&i.PickUpLongitude,
		&i.DropOffLatitude,
		&i.DropOffLongitude,
		&i.Status,
		&i.CreatedAt,
		&i.ExpireAt,
	)
	return i, err
}

const getNotificationSentByRequestID = `-- name: GetNotificationSentByRequestID :one
SELECT id, request_id, driver_phone_rejected FROM notification_sent
WHERE request_id = $1 LIMIT 1
`

func (q *Queries) GetNotificationSentByRequestID(ctx context.Context, requestID int64) (NotificationSent, error) {
	row := q.db.QueryRowContext(ctx, getNotificationSentByRequestID, requestID)
	var i NotificationSent
	err := row.Scan(&i.ID, &i.RequestID, pq.Array(&i.DriverPhoneRejected))
	return i, err
}

const updateNotificationSent = `-- name: UpdateNotificationSent :one
UPDATE notification_sent
SET driver_phone_rejected = $2
WHERE request_id = $1
RETURNING id, request_id, driver_phone_rejected
`

type UpdateNotificationSentParams struct {
	RequestID           int64    `json:"request_id"`
	DriverPhoneRejected []string `json:"driver_phone_rejected"`
}

func (q *Queries) UpdateNotificationSent(ctx context.Context, arg UpdateNotificationSentParams) (NotificationSent, error) {
	row := q.db.QueryRowContext(ctx, updateNotificationSent, arg.RequestID, pq.Array(arg.DriverPhoneRejected))
	var i NotificationSent
	err := row.Scan(&i.ID, &i.RequestID, pq.Array(&i.DriverPhoneRejected))
	return i, err
}
