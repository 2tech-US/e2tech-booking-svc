package services

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/client"
	"github.com/lntvan166/e2tech-booking-svc/internal/db"
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

	rsp, err := s.DriverSvc.UpdateDriverStatus(ctx, &client.UpdateDriverStatusRequest{
		Phone:  history.DriverPhone,
		Status: utils.DriverStatusFinding,
	})
	if err != nil {
		return &pb.CompleteTripResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to update driver status: %v", err),
		}, nil
	}
	if rsp.Status != http.StatusOK {
		return &pb.CompleteTripResponse{
			Status: rsp.Status,
			Error:  rsp.Error,
		}, nil
	}

	return &pb.CompleteTripResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) UpdateResponse(ctx context.Context, req *pb.UpdateResponseRequest) (*pb.UpdateResponseResponse, error) {
	response, err := s.DB.GetResponseByDriverPhone(ctx, req.DriverPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateResponseResponse{
				Status: http.StatusNotFound,
				Error:  "response not found",
			}, nil
		}
		return &pb.UpdateResponseResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get response: %v", err),
		}, nil
	}
	if response.DriverPhone != req.DriverPhone {
		return &pb.UpdateResponseResponse{
			Status: http.StatusBadRequest,
			Error:  "driver phone not match",
		}, nil
	}
	_, err = s.DB.UpdateResponse(ctx, db.UpdateResponseParams{
		ID:              response.ID,
		DriverPhone:     req.DriverPhone,
		DriverLatitude:  req.Latitude,
		DriverLongitude: req.Longitude,
	})
	if err != nil {
		return &pb.UpdateResponseResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to update response: %v", err),
		}, nil
	}

	return &pb.UpdateResponseResponse{
		Status: http.StatusOK,
	}, nil
}
