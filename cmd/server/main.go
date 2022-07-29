package main

import (
	"fmt"
	"log"
	"net"

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

	s := services.Server{
		DB:     DB,
		Config: &c,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBookingServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
