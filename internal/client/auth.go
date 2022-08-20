package client

import (
	"context"
	"fmt"

	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthServiceClient struct {
	Client pb.AuthServiceClient
}

func InitAuthServiceClient(c *config.Config) pb.AuthServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewAuthServiceClient(cc)
}

type GetDeviceTokenRequest struct {
	Phone string
}

func (s *AuthServiceClient) GetDeviceToken(ctx context.Context, req *GetDeviceTokenRequest) (*pb.GetDeviceTokenResponse, error) {
	return s.Client.GetDeviceToken(ctx, &pb.GetDeviceTokenRequest{
		Phone: req.Phone,
	})
}
