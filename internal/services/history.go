package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/db"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"github.com/lntvan166/e2tech-booking-svc/internal/utils"
	"go.uber.org/multierr"
)

func (s *Server) ListHistory(ctx context.Context, req *pb.ListHistoryRequest) (*pb.ListHistoryResponse, error) {
	histories := []db.History{}
	startDate, err1 := utils.ParseStringToDate(req.StartDate)
	endDate, err2 := utils.ParseStringToDate(req.EndDate)
	err := multierr.Combine(
		err1, err2,
	)

	if err != nil {
		return &pb.ListHistoryResponse{
			Status: http.StatusBadRequest,
			Error:  fmt.Sprintf("invalid date: %v", err),
		}, nil
	}
	if req.Role == utils.PASSENGER {
		var err error
		histories, err = s.DB.ListHistoryByPassengerPhone(ctx, db.ListHistoryByPassengerPhoneParams{
			PassengerPhone: req.Phone,
			StartDate:      startDate,
			EndDate:        endDate,
			Limit:          req.Limit,
			Offset:         req.Offset,
		})
		if err != nil {
			return &pb.ListHistoryResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("failed to list history: %v", err),
			}, nil
		}
	}
	if req.Role == utils.DRIVER {
		var err error
		histories, err = s.DB.ListHistoryByDriverPhone(ctx, db.ListHistoryByDriverPhoneParams{
			DriverPhone: req.Phone,
			StartDate:   startDate,
			EndDate:     endDate,
			Limit:       req.Limit,
			Offset:      req.Offset,
		})
		if err != nil {
			return &pb.ListHistoryResponse{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("failed to list history: %v", err),
			}, nil
		}
	}

	dataRsp := make([]*pb.History, len(histories))
	for i, h := range histories {
		dataRsp[i] = &pb.History{
			Type:           h.Type,
			PassengerPhone: h.PassengerPhone,
			DriverPhone:    h.DriverPhone,
			PickUpLocation: &pb.Location{
				Latitude:  h.PickUpLatitude,
				Longitude: h.PickUpLongitude,
			},
			DropOffLocation: &pb.Location{
				Latitude:  h.DropOffLatitude,
				Longitude: h.DropOffLongitude,
			},
			CreatedAt: utils.ParsedDateToString(h.CreatedAt),
			DoneAt:    utils.ParsedDateToString(h.DoneAt),
		}
	}

	return &pb.ListHistoryResponse{
		Status:  http.StatusOK,
		History: dataRsp,
	}, nil
}
