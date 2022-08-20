package client

import (
	"fmt"

	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PassengerServiceClient struct {
	Client pb.PassengerServiceClient
}

func InitPassengerServiceClient(c *config.Config) pb.PassengerServiceClient {
	cc, err := grpc.Dial(c.PassengerSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewPassengerServiceClient(cc)
}
