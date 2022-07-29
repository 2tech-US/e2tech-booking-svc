// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: response.sql

package db

import (
	"context"
)

const createResponse = `-- name: CreateResponse :one
INSERT INTO response (
  request_id,
  driver_id,
  driver_latitude,
  driver_longitude
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at
`

type CreateResponseParams struct {
	RequestID       int64   `json:"request_id"`
	DriverID        int64   `json:"driver_id"`
	DriverLatitude  float64 `json:"driver_latitude"`
	DriverLongitude float64 `json:"driver_longitude"`
}

func (q *Queries) CreateResponse(ctx context.Context, arg CreateResponseParams) (Response, error) {
	row := q.db.QueryRowContext(ctx, createResponse,
		arg.RequestID,
		arg.DriverID,
		arg.DriverLatitude,
		arg.DriverLongitude,
	)
	var i Response
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		&i.DriverID,
		&i.DriverName,
		&i.DriverLatitude,
		&i.DriverLongitude,
		&i.CreatedAt,
	)
	return i, err
}

const deleteResponse = `-- name: DeleteResponse :exec
DELETE FROM response
WHERE id = $1
`

func (q *Queries) DeleteResponse(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteResponse, id)
	return err
}

const getResponse = `-- name: GetResponse :one
SELECT id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at FROM response
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetResponse(ctx context.Context, id int64) (Response, error) {
	row := q.db.QueryRowContext(ctx, getResponse, id)
	var i Response
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		&i.DriverID,
		&i.DriverName,
		&i.DriverLatitude,
		&i.DriverLongitude,
		&i.CreatedAt,
	)
	return i, err
}

const getResponseByPassengerID = `-- name: GetResponseByPassengerID :one
SELECT id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at FROM response
WHERE request_id = $1 LIMIT 1
`

func (q *Queries) GetResponseByPassengerID(ctx context.Context, requestID int64) (Response, error) {
	row := q.db.QueryRowContext(ctx, getResponseByPassengerID, requestID)
	var i Response
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		&i.DriverID,
		&i.DriverName,
		&i.DriverLatitude,
		&i.DriverLongitude,
		&i.CreatedAt,
	)
	return i, err
}

const getResponseByRequestID = `-- name: GetResponseByRequestID :one
SELECT id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at FROM response
WHERE request_id = $1 LIMIT 1
`

func (q *Queries) GetResponseByRequestID(ctx context.Context, requestID int64) (Response, error) {
	row := q.db.QueryRowContext(ctx, getResponseByRequestID, requestID)
	var i Response
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		&i.DriverID,
		&i.DriverName,
		&i.DriverLatitude,
		&i.DriverLongitude,
		&i.CreatedAt,
	)
	return i, err
}

const listResponses = `-- name: ListResponses :many
SELECT id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at FROM response
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListResponsesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListResponses(ctx context.Context, arg ListResponsesParams) ([]Response, error) {
	rows, err := q.db.QueryContext(ctx, listResponses, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Response
	for rows.Next() {
		var i Response
		if err := rows.Scan(
			&i.ID,
			&i.RequestID,
			&i.DriverID,
			&i.DriverName,
			&i.DriverLatitude,
			&i.DriverLongitude,
			&i.CreatedAt,
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

const updateResponse = `-- name: UpdateResponse :one
UPDATE response
SET driver_id = $2,
    driver_latitude = $3,
    driver_longitude = $4
WHERE id = $1
RETURNING id, request_id, driver_id, driver_name, driver_latitude, driver_longitude, created_at
`

type UpdateResponseParams struct {
	ID              int64   `json:"id"`
	DriverID        int64   `json:"driver_id"`
	DriverLatitude  float64 `json:"driver_latitude"`
	DriverLongitude float64 `json:"driver_longitude"`
}

func (q *Queries) UpdateResponse(ctx context.Context, arg UpdateResponseParams) (Response, error) {
	row := q.db.QueryRowContext(ctx, updateResponse,
		arg.ID,
		arg.DriverID,
		arg.DriverLatitude,
		arg.DriverLongitude,
	)
	var i Response
	err := row.Scan(
		&i.ID,
		&i.RequestID,
		&i.DriverID,
		&i.DriverName,
		&i.DriverLatitude,
		&i.DriverLongitude,
		&i.CreatedAt,
	)
	return i, err
}