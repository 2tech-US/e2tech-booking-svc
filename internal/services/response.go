package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"github.com/lntvan166/e2tech-booking-svc/internal/utils"
)

func (s *Server) GetResponse(ctx context.Context, req *pb.GetResponseRequest) (*pb.GetResponseResponse, error) {
	request, err := s.DB.GetRequestByPassengerID(ctx, utils.NullInt64(req.PassengerId))
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetResponseResponse{
				Status: http.StatusNotFound,
				Error:  "request not found",
			}, nil
		}
		return &pb.GetResponseResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get request: %v", err),
		}, nil
	}

	response, err := s.DB.GetResponseByRequestID(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetResponseResponse{
				Status: http.StatusNotFound,
				Error:  "response not found",
			}, nil
		}
		return &pb.GetResponseResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get response: %v", err),
		}, nil
	}

	dataRsp := &pb.DriverResponse{
		Status: request.Status,
		Name:   response.DriverName,
		Location: &pb.Location{
			Latitude:  response.DriverLatitude,
			Longitude: response.DriverLongitude,
		},
	}

	return &pb.GetResponseResponse{
		Status: http.StatusOK,
		Driver: dataRsp,
	}, nil
}
