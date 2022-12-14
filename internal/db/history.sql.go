// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: history.sql

package db

import (
	"context"
	"time"
)

const createHistory = `-- name: CreateHistory :one
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
RETURNING id, type, passenger_phone, driver_phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, price, created_at, done_at
`

type CreateHistoryParams struct {
	Type             string    `json:"type"`
	PassengerPhone   string    `json:"passenger_phone"`
	DriverPhone      string    `json:"driver_phone"`
	PickUpLatitude   float64   `json:"pick_up_latitude"`
	PickUpLongitude  float64   `json:"pick_up_longitude"`
	DropOffLatitude  float64   `json:"drop_off_latitude"`
	DropOffLongitude float64   `json:"drop_off_longitude"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
	DoneAt           time.Time `json:"done_at"`
}

func (q *Queries) CreateHistory(ctx context.Context, arg CreateHistoryParams) (History, error) {
	row := q.db.QueryRowContext(ctx, createHistory,
		arg.Type,
		arg.PassengerPhone,
		arg.DriverPhone,
		arg.PickUpLatitude,
		arg.PickUpLongitude,
		arg.DropOffLatitude,
		arg.DropOffLongitude,
		arg.Price,
		arg.CreatedAt,
		arg.DoneAt,
	)
	var i History
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.PassengerPhone,
		&i.DriverPhone,
		&i.PickUpLatitude,
		&i.PickUpLongitude,
		&i.DropOffLatitude,
		&i.DropOffLongitude,
		&i.Price,
		&i.CreatedAt,
		&i.DoneAt,
	)
	return i, err
}

const deleteHistory = `-- name: DeleteHistory :exec
DELETE FROM history
WHERE id = $1
`

func (q *Queries) DeleteHistory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteHistory, id)
	return err
}

const getHistory = `-- name: GetHistory :one
SELECT id, type, passenger_phone, driver_phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, price, created_at, done_at
FROM history
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetHistory(ctx context.Context, id int64) (History, error) {
	row := q.db.QueryRowContext(ctx, getHistory, id)
	var i History
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.PassengerPhone,
		&i.DriverPhone,
		&i.PickUpLatitude,
		&i.PickUpLongitude,
		&i.DropOffLatitude,
		&i.DropOffLongitude,
		&i.Price,
		&i.CreatedAt,
		&i.DoneAt,
	)
	return i, err
}

const listHistories = `-- name: ListHistories :many
SELECT id, type, passenger_phone, driver_phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, price, created_at, done_at
FROM history
WHERE created_at >= $3
  AND created_at <= $4
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListHistoriesParams struct {
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (q *Queries) ListHistories(ctx context.Context, arg ListHistoriesParams) ([]History, error) {
	rows, err := q.db.QueryContext(ctx, listHistories,
		arg.Limit,
		arg.Offset,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.PassengerPhone,
			&i.DriverPhone,
			&i.PickUpLatitude,
			&i.PickUpLongitude,
			&i.DropOffLatitude,
			&i.DropOffLongitude,
			&i.Price,
			&i.CreatedAt,
			&i.DoneAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listHistoryByDriverPhone = `-- name: ListHistoryByDriverPhone :many
SELECT id, type, passenger_phone, driver_phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, price, created_at, done_at
FROM history
WHERE driver_phone = $1
  AND created_at >= $4
  AND created_at <= $5
ORDER BY id
LIMIT $2 OFFSET $3
`

type ListHistoryByDriverPhoneParams struct {
	DriverPhone string    `json:"driver_phone"`
	Limit       int32     `json:"limit"`
	Offset      int32     `json:"offset"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

func (q *Queries) ListHistoryByDriverPhone(ctx context.Context, arg ListHistoryByDriverPhoneParams) ([]History, error) {
	rows, err := q.db.QueryContext(ctx, listHistoryByDriverPhone,
		arg.DriverPhone,
		arg.Limit,
		arg.Offset,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.PassengerPhone,
			&i.DriverPhone,
			&i.PickUpLatitude,
			&i.PickUpLongitude,
			&i.DropOffLatitude,
			&i.DropOffLongitude,
			&i.Price,
			&i.CreatedAt,
			&i.DoneAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listHistoryByPassengerPhone = `-- name: ListHistoryByPassengerPhone :many
SELECT id, type, passenger_phone, driver_phone, pick_up_latitude, pick_up_longitude, drop_off_latitude, drop_off_longitude, price, created_at, done_at
FROM history
WHERE passenger_phone = $1
  AND created_at >= $4
  AND created_at <= $5
ORDER BY id
LIMIT $2 OFFSET $3
`

type ListHistoryByPassengerPhoneParams struct {
	PassengerPhone string    `json:"passenger_phone"`
	Limit          int32     `json:"limit"`
	Offset         int32     `json:"offset"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}

func (q *Queries) ListHistoryByPassengerPhone(ctx context.Context, arg ListHistoryByPassengerPhoneParams) ([]History, error) {
	rows, err := q.db.QueryContext(ctx, listHistoryByPassengerPhone,
		arg.PassengerPhone,
		arg.Limit,
		arg.Offset,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []History
	for rows.Next() {
		var i History
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.PassengerPhone,
			&i.DriverPhone,
			&i.PickUpLatitude,
			&i.PickUpLongitude,
			&i.DropOffLatitude,
			&i.DropOffLongitude,
			&i.Price,
			&i.CreatedAt,
			&i.DoneAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
