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

	if len(driverRsp.Drivers) == 0 {
		return &pb.CreateRequestResponse{
			Status: http.StatusNotFound,
			Error:  "no driver found",
		}, nil
	}

	arg := db.CreateRequestParams{
		Type:             req.Type,
		Phone:            req.Phone,
		PickUpLatitude:   req.PickUpLocation.Latitude,
		PickUpLongitude:  req.PickUpLocation.Longitude,
		DropOffLatitude:  req.DropOffLocation.Latitude,
		DropOffLongitude: req.DropOffLocation.Longitude,
	}

	_, err = s.DB.CreateRequest(ctx, arg)
	if err != nil {
		return &pb.CreateRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to create request: %v", err),
		}, nil
	}
	for _, driver := range driverRsp.Drivers {
		status, err := s.sendNotification(ctx, driver.Phone, sendNotificationData{
			Title: "A new passenger come to pick you up",
			Body:  fmt.Sprintf("Passenger %v away from you", driver.Distance),
			Data: map[string]interface{}{
				"passenger_phone": driver.Phone,
				"distance":        driver.Distance,
				"pickup_lat":      req.PickUpLocation.Latitude,
				"pickup_lng":      req.PickUpLocation.Longitude,
				"dropoff_lat":     req.DropOffLocation.Latitude,
				"dropoff_lng":     req.DropOffLocation.Longitude,
			},
		})
		if err != nil {
			return &pb.CreateRequestResponse{
				Status: status,
				Error:  err.Error(),
			}, nil
		}
	}

	return &pb.CreateRequestResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) CloseRequest(ctx context.Context, req *pb.CloseRequestRequest) (*pb.CloseRequestResponse, error) {
	err := s.DB.CloseRequest(ctx, req.Phone)
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

	_, err = s.DriverSvc.UpdateDriverStatus(ctx, &client.UpdateDriverStatusRequest{
		Phone:  req.DriverPhone,
		Status: utils.DriverStatusInProgress,
	})
	if err != nil {
		return &pb.AcceptRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to update driver status: %v", err),
		}, nil
	}

	if request.Type == utils.RequestTypeApp {
		status, err := s.sendNotification(ctx, req.PassengerPhone, sendNotificationData{
			Title: "Your driver is on the way",
			Body:  "Please wait for a while",
			Data: map[string]interface{}{
				"driver_name": driver.Driver.Name,
				"driver_phone":// Used to update the notification sent table.
				req.DriverPhone,
				"driver_lat": driverLocation.Latitude,
				"driver_lng": driverLocation.Longitude,
			},
		})
		if err != nil {
			return &pb.AcceptRequestResponse{
				Status: status,
				Error:  fmt.Sprintf("failed to send notification: %v", err),
			}, nil
		}
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

	status, err := s.updateNotificationRejected(ctx, request.ID, req.DriverPhone)
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

func (s *Server) updateNotificationRejected(ctx context.Context, requestID int64, driverPhone string) (int64, error) {
	var driversRejected []string
	notificationRejected, err := s.DB.GetNotificationSentByRequestID(ctx, requestID)
	if err != nil {
		if err != sql.ErrNoRows {
			return http.StatusInternalServerError, fmt.Errorf("failed to get notification sent: %v", err)
		}
		driversRejected = []string{driverPhone}
		_, err = s.DB.CreateNotificationSent(ctx, db.CreateNotificationSentParams{
			RequestID:           requestID,
			DriverPhoneRejected: driversRejected,
		})
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to create notification sent: %v", err)
		}
		return http.StatusOK, nil
	}

	driversRejected = notificationRejected.DriverPhoneRejected
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

func (s *Server) GetRequest(ctx context.Context, req *pb.GetRequestRequest) (*pb.GetRequestResponse, error) {
	request, err := s.DB.GetRequestByPhone(ctx, req.PassengerPhone)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetRequestResponse{
				Status: http.StatusNotFound,
				Error:  "request not found: %v",
			}, nil
		}
		return &pb.GetRequestResponse{
			Status: http.StatusInternalServerError,
			Error:  fmt.Sprintf("failed to get request: %v", err),
		}, nil
	}

	return &pb.GetRequestResponse{
		Status: http.StatusOK,
		Request: &pb.Request{
			Id:    request.ID,
			Type:  request.Type,
			Phone: request.Phone,
			PickUpLocation: &pb.Location{
				Latitude:  request.PickUpLatitude,
				Longitude: request.PickUpLongitude,
			},
			DropOffLocation: &pb.Location{
				Latitude:  request.DropOffLatitude,
				Longitude: request.DropOffLongitude,
			},
			CreatedAt: utils.ParsedDateToString(request.CreatedAt),
		},
	}, nil
}
