package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lntvan166/e2tech-booking-svc/internal/client"
	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/lntvan166/e2tech-booking-svc/internal/db"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"github.com/lntvan166/e2tech-booking-svc/internal/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	DB := db.Connect(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Booking Svc on", c.Port)

	authSvc := &client.AuthServiceClient{
		Client: client.InitAuthServiceClient(&c),
	}
	passengerSvc := &client.PassengerServiceClient{
		Client: client.InitPassengerServiceClient(&c),
	}
	driverSvc := &client.DriverServiceClient{
		Client: client.InitDriverServiceClient(&c),
	}
	notificationSvc := &client.NotificationServiceClient{
		Client: client.InitNotificationServiceClient(&c),
		ApiKey: c.FirebaseApiKey,
	}

	s := services.Server{
		DB:              DB,
		AuthSvc:         authSvc,
		PassengerSvc:    passengerSvc,
		DriverSvc:       driverSvc,
		NotificationSvc: notificationSvc,
		Config:          &c,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBookingServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
