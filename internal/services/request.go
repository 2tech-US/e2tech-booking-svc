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

func (s *Server) CreateRequest(ctx context.Context, req *pb.CreateRequestRequest) (*pb.CreateRequestResponse, error) {
	arg := db.CreateRequestParams{
		Type:             req.Type,
		Phone:            req.Phone,
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

	driverRsp, err := s.DriverSvc.GetDriverNearby(ctx, &client.GetDriverNearbyRequest{
		Latitude:        req.PickUpLocation.Latitude,
		Longitude:       req.PickUpLocation.Longitude,
		NumberOfDrivers: 5,
	})
	if err != nil {
		return &pb.CreateRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get driver nearby: %v", err),
		}, nil
	}

	status, err := s.sendNotification(ctx, driverRsp.Drivers, sendNotificationData{
		Phone:            req.Phone,
		PickUpLatitude:   req.PickUpLocation.Latitude,
		PickUpLongitude:  req.PickUpLocation.Longitude,
		DropOffLatitude:  req.DropOffLocation.Latitude,
		DropOffLongitude: req.DropOffLocation.Longitude,
	})
	if err != nil {
		return &pb.CreateRequestResponse{
			Status: status,
			Error:  err.Error(),
		}, nil
	}

	return &pb.CreateRequestResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) CloseRequest(ctx context.Context, req *pb.CloseRequestRequest) (*pb.CloseRequestResponse, error) {
	request, err := s.DB.GetRequestByPhone(ctx, req.Phone)
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

	err = s.DB.DeleteRequest(ctx, request.Phone)
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
func (s *Server) AcceptRequest(ctx context.Context, req *pb.AcceptRequestRequest) (*pb.AcceptRequestResponse, error) {
	request, err := s.DB.GetRequestByPhone(ctx, req.PassengerPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.AcceptRequestResponse{
				Status: http.StatusNotFound,
				Error:  "request not found: %v",
			}, nil
		}
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get request: %v", err),
		}, nil
	}

	if request.Status != utils.RequestStatusFinding {
		return &pb.AcceptRequestResponse{
			Status: http.StatusBadRequest,
			Error:  "request is accepted or closed",
		}, nil
	}

	driverLocation, err := s.DriverSvc.GetLocation(ctx, &client.GetLocationRequest{
		Phone: req.DriverPhone,
	})
	if err != nil {
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get driver location: %v", err),
		}, nil
	}

	driver, err := s.DriverSvc.GetDriver(ctx, &client.GetDriverRequest{
		Phone: req.DriverPhone,
	})
	if err != nil {
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get driver: %v", err),
		}, nil
	}

	arg := db.CreateResponseParams{
		RequestID:       request.ID,
		DriverName:      driver.Driver.Name,
		DriverPhone:     req.DriverPhone,
		DriverLatitude:  driverLocation.Latitude,
		DriverLongitude: driverLocation.Longitude,
	}

	_, err = s.DB.CreateResponse(ctx, arg)
	if err != nil {
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to accept request: %v", err),
		}, nil
	}

	_, err = s.DB.UpdateStatusRequest(ctx, db.UpdateStatusRequestParams{
		Phone:  req.PassengerPhone,
		Status: utils.DriverStatusInProgress,
	})
	if err != nil {
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to update request status: %v", err),
		}, nil
	}

	return &pb.AcceptRequestResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) RejectRequest(ctx context.Context, req *pb.RejectRequestRequest) (*pb.RejectRequestResponse, error) {
	request, err := s.DB.GetRequestByPhone(ctx, req.PassengerPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.RejectRequestResponse{
				Status: http.StatusNotFound,
				Error:  "request is accepted or closed: %v",
			}, nil
		}
		return &pb.RejectRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get request: %v", err),
		}, nil
	}

	if request.Status != utils.RequestStatusFinding {
		return &pb.RejectRequestResponse{
			Status: http.StatusBadRequest,
			Error:  "request is accepted or closed",
		}, nil
	}

	status, err := s.updateNotificationRejected(ctx, req.PassengerPhone, req.DriverPhone)
	if err != nil {
		return &pb.RejectRequestResponse{
			Status: status,
			Error:  err.Error(),
		}, nil
	}

	status, err = s.ResendNotification(ctx, request.Phone, req.DriverPhone)
	if err != nil {
		return &pb.RejectRequestResponse{
			Status: status,
			Error:  err.Error(),
		}, nil
	}

	return &pb.RejectRequestResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) updateNotificationRejected(ctx context.Context, passengerPhone, driverPhone string) (int64, error) {
	notificationRejected, err := s.DB.GetNotificationSentByPassengerPhone(ctx, passengerPhone)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get notification sent: %v", err)
	}

	driversRejected := notificationRejected.DriverPhoneRejected
	driversRejected = append(driversRejected, driverPhone)

	_, err = s.DB.UpdateNotificationSent(ctx, db.UpdateNotificationSentParams{
		RequestID:           notificationRejected.RequestID,
		DriverPhoneRejected: driversRejected,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to update notification sent: %v", err)
	}

	return http.StatusOK, nil
}
