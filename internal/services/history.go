package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"github.com/lntvan166/e2tech-booking-svc/internal/utils"
)

func (s *Server) ListHistory(ctx context.Context, req *pb.ListHistoryRequest) (*pb.ListHistoryResponse, error) {
	histories, err := s.DB.ListHistoryByPassengerID(ctx, req.PassengerId)
	if err != nil {
		return &pb.ListHistoryResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to list history: %v", err),
		}, nil
	}

	dataRsp := make([]*pb.History, len(histories))
	for i, h := range histories {
		dataRsp[i] = &pb.History{
			Type:        h.Type,
			PassengerId: h.PassengerID,
			Phone:       h.Phone,
			PickUpLocation: &pb.Location{
				Latitude:  h.PickUpLatitude,
				Longitude: h.PickUpLongitude,
			},
			DropOffLocation: &pb.Location{
				Latitude:  h.DropOffLatitude,
				Longitude: h.DropOffLongitude,
			},
			DriverId:  h.DriverID,
			CreatedAt: utils.ParsedDateToString(h.CreatedAt),
		}
	}

	return &pb.ListHistoryResponse{
		Status:  http.StatusOK,
		History: dataRsp,
	}, nil
}
