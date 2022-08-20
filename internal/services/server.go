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
	NotificationSvc *client.NotificationServiceClient
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

	status, err := s.sendNotification(ctx, driverNearby, sendNotificationData{
		Phone:            passengerPhone,
		PickUpLatitude:   request.PickUpLatitude,
		PickUpLongitude:  request.PickUpLongitude,
		DropOffLatitude:  request.DropOffLatitude,
		DropOffLongitude: request.DropOffLongitude,
	})
	if err != nil {
		return status, fmt.Errorf("failed to send notification: %v", err)
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
	Phone            string
	PickUpLatitude   float64
	PickUpLongitude  float64
	DropOffLatitude  float64
	DropOffLongitude float64
}

func (s *Server) sendNotification(ctx context.Context, drivers []*pb.DriverNearby, data sendNotificationData) (int64, error) {
	for _, driver := range drivers {
		authRsp, err := s.AuthSvc.GetDeviceToken(ctx, &client.GetDeviceTokenRequest{
			Phone: driver.Phone,
		})
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to get device token: %v", err)
		}
		if authRsp.Token == "" {
			return http.StatusInternalServerError, fmt.Errorf("failed to get device token: token is empty")
		}

		// send notification to driver
		_, err = s.NotificationSvc.SendNotification(ctx, &client.SendNotificationRequest{
			Tokens:  []string{authRsp.Token},
			Message: "You have a new passenger",
			Data: map[string]interface{}{
				"passenger_phone": data.Phone,
				"distance":        driver.Distance,
				"pickup_lat":      data.PickUpLatitude,
				"pickup_lng":      data.PickUpLongitude,
				"dropoff_lat":     data.DropOffLatitude,
				"dropoff_lng":     data.DropOffLongitude,
			},
		})
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to send notification: %v", err)
		}
	}
	return http.StatusOK, nil
}
