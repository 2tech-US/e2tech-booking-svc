package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lntvan166/e2tech-booking-svc/internal/utils"
)

func (q *Queries) MakeHistory(ctx context.Context, passengerPhone, driverPhone string) (History, error) {
	tx, err := db.Begin()
	if err != nil {
		return History{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	qtx := q.WithTx(tx)

	request, err := qtx.GetRequestByPhone(ctx, passengerPhone)
	if err != nil {
		return History{}, err
	}

	response, err := qtx.GetResponseByRequestID(ctx, request.ID)
	if err != nil {
		return History{}, err
	}

	if response.DriverPhone != driverPhone {
		return History{}, fmt.Errorf("driver phone does not match")
	}

	history, err := qtx.CreateHistory(ctx, CreateHistoryParams{
		Type:             request.Type,
		PassengerPhone:   request.Phone,
		DriverPhone:      response.DriverPhone,
		PickUpLatitude:   request.PickUpLatitude,
		PickUpLongitude:  request.PickUpLongitude,
		DropOffLatitude:  request.DropOffLatitude,
		DropOffLongitude: request.DropOffLongitude,
		Price:            utils.CalculatePrice(request.PickUpLatitude, request.PickUpLongitude, request.DropOffLatitude, request.DropOffLongitude),
		CreatedAt:        request.CreatedAt,
		DoneAt:           time.Now(),
	})
	if err != nil {
		return History{}, err
	}
	err = qtx.DeleteResponse(ctx, response.DriverPhone)
	if err != nil {
		return History{}, err
	}
	err = qtx.DeleteRequest(ctx, request.Phone)
	if err != nil {
		return History{}, err
	}
	err = tx.Commit()
	if err != nil {
		return History{}, err
	}
	return history, nil

}

func (q *Queries) CloseRequest(ctx context.Context, phone string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	qtx := q.WithTx(tx)

	request, err := qtx.GetRequestByPhone(ctx, phone)
	if err != nil {
		return err
	}

	err = qtx.DeleteResponseByRequestID(ctx, request.ID)
	if err != nil {
		return err
	}
	err = qtx.DeleteRequest(ctx, phone)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
