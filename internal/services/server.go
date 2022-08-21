package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lntvan166/e2tech-booking-svc/internal/client"
	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/lntvan166/e2tech-booking-svc/internal/db"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
)

type Server struct {
	DB              *db.Queries
	Config          *config.Config
	AuthSvc         *client.AuthServiceClient
	PassengerSvc    *client.PassengerServiceClient
	DriverSvc       *client.DriverServiceClient
	NotificationSvc *client.NotificationServiceClientV2
	pb.UnimplementedBookingServiceServer
}

func (s *Server) ResendNotification(ctx context.Context, passengerPhone, driverPhone string) (int64, error) {
	request, err := s.DB.GetRequestByPhone(ctx, passengerPhone)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get request: %v", err)
	}

	driverRsp, err := s.DriverSvc.GetDriverNearby(ctx, &client.GetDriverNearbyRequest{
		Latitude:        request.PickUpLatitude,
		Longitude:       request.PickUpLongitude,
		NumberOfDrivers: 20,
	})
	if err != nil {
		return http.StatusBadGateway, fmt.Errorf("failed to get driver nearby: %v", err)
	}
	driverNearby := driverRsp.Drivers
	notificationRejected, err := s.DB.GetNotificationSentByPassengerPhone(ctx, passengerPhone)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get notification sent: %v", err)
	}
	for _, driver := range notificationRejected.DriverPhoneRejected {
		driverNearby = removeDriver(driverNearby, driver)
	}

	// limit by 5
	if len(driverNearby) > 5 {
		driverNearby = driverNearby[:5]
	}

	for _, driver := range driverNearby {
		status, err := s.sendNotification(ctx, driverPhone, sendNotificationData{
			Title: "A new passenger come to pick you up",
			Body:  fmt.Sprintf("Passenger %v away from you", driver.Distance),
			Data: map[string]interface{}{
				"passenger_phone": driver.Phone,
				"distance":        driver.Distance,
				"pickup_lat":      request.PickUpLatitude,
				"pickup_lng":      request.PickUpLongitude,
				"dropoff_lat":     request.DropOffLatitude,
				"dropoff_lng":     request.DropOffLongitude,
			},
		})
		if err != nil {
			return status, fmt.Errorf("failed to send notification: %v", err)
		}
	}

	return http.StatusOK, nil
}

func removeDriver(drivers []*pb.DriverNearby, driver string) []*pb.DriverNearby {
	for i, d := range drivers {
		if d.Phone == driver {
			drivers = append(drivers[:i], drivers[i+1:]...)
			break
		}
	}
	return drivers
}

type sendNotificationData struct {
	Title string
	Body  string

	Data map[string]interface{}
}

func (s *Server) sendNotification(ctx context.Context, Phone string, req sendNotificationData) (int64, error) {
	authRsp, err := s.AuthSvc.GetDeviceToken(ctx, &client.GetDeviceTokenRequest{
		Phone: Phone,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to get device token: %v", err)
	}
	if len(authRsp.Token) < 32 {
		return http.StatusOK, nil
	}

	// send notification to driver
	rsp, err := s.NotificationSvc.SendNotificationV2(ctx, &client.SendNotificationRequestV2{
		To:    authRsp.Token,
		Title: req.Title,
		Body:  req.Body,
		Data:  req.Data,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to send notification: %v", err)
	}

	if rsp.Failure == 1 {
		return http.StatusInternalServerError, fmt.Errorf("failed to send notification: %v", rsp.Error)
	}

	return http.StatusOK, nil
}
