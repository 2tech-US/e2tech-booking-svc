package client

import (
	"context"
	"fmt"

	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DriverServiceClient struct {
	Client pb.DriverServiceClient
}

func InitDriverServiceClient(c *config.Config) pb.DriverServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.DriverSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewDriverServiceClient(cc)
}

type GetDriverNearbyRequest struct {
	Latitude        float64
	Longitude       float64
	NumberOfDrivers int32
}

func (s *DriverServiceClient) GetDriverNearby(ctx context.Context, req *GetDriverNearbyRequest) (*pb.GetDriverNearbyResponse, error) {
	return s.Client.GetDriverNearby(ctx, &pb.GetDriverNearbyRequest{
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		NumberOfDrivers: req.NumberOfDrivers,
	})
}

type GetDriverRequest struct {
	Phone string
}

func (s *DriverServiceClient) GetDriver(ctx context.Context, req *GetDriverRequest) (*pb.GetDriverByPhoneResponse, error) {
	return s.Client.GetDriverByPhone(ctx, &pb.GetDriverByPhoneRequest{
		Phone: req.Phone,
	})
}

type GetLocationRequest struct {
	Phone string
}

func (s *DriverServiceClient) GetLocation(ctx context.Context, req *GetLocationRequest) (*pb.GetLocationResponse, error) {
	return s.Client.GetLocation(ctx, &pb.GetLocationRequest{
		Phone: req.Phone,
	})
}
