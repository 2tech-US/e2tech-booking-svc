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
	request, err := s.DB.GetRequestByPhone(ctx, req.Phone)
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

func (s *Server) CompleteTrip(ctx context.Context, req *pb.CompleteTripRequest) (*pb.CompleteTripResponse, error) {
	history, err := s.DB.MakeHistory(ctx, req.PassengerPhone, req.DriverPhone)
	if err != nil {
		return &pb.CompleteTripResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to make history: %v", err),
		}, nil
	}

	if history.Type == utils.RequestTypeApp {
		s.sendNotification(ctx, history.PassengerPhone, sendNotificationData{
			Title: "Trip Completed",
			Body:  "Your trip has been completed",
			Data: map[string]interface{}{
				"price":        history.Price,
				"driver_phone": history.DriverPhone,
			},
		})
	}

	return &pb.CompleteTripResponse{
		Status: http.StatusOK,
	}, nil
}
