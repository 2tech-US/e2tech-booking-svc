package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/db"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"github.com/lntvan166/e2tech-booking-svc/internal/utils"
)

func (s *Server) CreateRequest(ctx context.Context, req *pb.CreateRequestRequest) (*pb.CreateRequestResponse, error) {
	arg := db.CreateRequestParams{
		Type:             req.Type,
		PassengerID:      utils.NullInt64(req.PassengerId),
		PickUpLatitude:   req.PickUpLocation.Latitude,
		PickUpLongitude:  req.PickUpLocation.Longitude,
		DropOffLatitude:  req.DropOffLocation.Latitude,
		DropOffLongitude: req.DropOffLocation.Longitude,
	}

	_, err := s.DB.CreateRequest(ctx, arg)
	if err != nil {
		return &pb.CreateRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to create request: %v", err),
		}, nil
	}

	return &pb.CreateRequestResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) CloseRequest(ctx context.Context, req *pb.CloseRequestRequest) (*pb.CloseRequestResponse, error) {
	request, err := s.DB.GetRequestByPassengerID(ctx, utils.NullInt64(req.PassengerId))
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.CloseRequestResponse{
				Status: http.StatusNotFound,
				Error:  "request not found: %v",
			}, nil
		}
		return &pb.CloseRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get request: %v", err),
		}, nil
	}

	err = s.DB.DeleteRequest(ctx, request.ID)
	if err != nil {
		return &pb.CloseRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to close request: %v", err),
		}, nil
	}

	return &pb.CloseRequestResponse{
		Status: http.StatusOK,
	}, nil
}

// TODO driver service
// func (s *Server) AcceptRequest(ctx context.Context, req *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
// 	request, err := s.DB.GetRequest(ctx, req.RequestId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return &pb.AcceptRequestResponse{
// 				Status: http.StatusNotFound,
// 				Error:  "request not found: %v",
// 			}, nil
// 		}
// 		return &pb.AcceptRequestResponse{
// 			Status: http.StatusInternalServerError,
// 			Error:  fmt.Sprintf("failed to get request: %v", err),
// 		}, nil
// 	}
// 	arg := db.CreateResponseParams{
// 		RequestID:      request.ID,
// 		DriverID:       req.DriverId,
// 		DriverLatitude: req.DriverLocation.Latitude,
// 	}

// 	_, err = s.DB.CreateResponse(ctx, arg)
// 	if err != nil {
// 		return &pb.AcceptRequestResponse{
// 			Status: http.StatusInternalServerError,
// 			Error:  fmt.Sprintf("failed to accept request: %v", err),
// 		}, nil
// 	}

// 	return &pb.AcceptRequestResponse{
// 		Status: http.StatusOK,
// 	}, nil
// }
