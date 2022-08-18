package client

import (
	"context"

	"github.com/lntvan166/e2tech-booking-svc/internal/pb"
)

type GetDriverNearbyRequest struct {
	Latitude  float64
	Longitude float64
}

func (s *DriverServiceClient) GetDriverNearby(ctx context.Context, req *GetDriverNearbyRequest) (*pb.GetDriverNearbyResponse, error) {
	return s.Client.GetDriverNearby(ctx, &pb.GetDriverNearbyRequest{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
}
