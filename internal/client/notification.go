package client

import (
	"context"
	"fmt"

	"github.com/appleboy/gorush/rpc/proto"
	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

type NotificationServiceClient struct {
	Client proto.GorushClient
	ApiKey string
}

func InitNotificationServiceClient(c *config.Config) proto.GorushClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.GorushUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return proto.NewGorushClient(cc)
}

type SendNotificationRequest struct {
	Tokens  []string
	Message string
	Data    map[string]interface{}
}

func (s *NotificationServiceClient) SendNotification(ctx context.Context, req *SendNotificationRequest) (*proto.NotificationReply, error) {
	var data *structpb.Struct
	if len(req.Data) > 0 {
		data = &structpb.Struct{
			Fields: make(map[string]*structpb.Value),
		}
		for k, v := range req.Data {
			data.Fields[k] = &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: v.(string),
				},
			}
		}
	}

	return s.Client.Send(ctx, &proto.NotificationRequest{
		Platform: 2,
		Tokens:   req.Tokens,
		Title:    "A new passenger",
		Message:  req.Message,
		Priority: proto.NotificationRequest_HIGH,
		Data:     data,
	})
}
