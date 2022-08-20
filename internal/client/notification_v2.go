package client

import (
	"context"

	"github.com/appleboy/go-fcm"
	"github.com/lntvan166/e2tech-booking-svc/internal/config"
)

type NotificationServiceClientV2 struct {
	Client fcm.Client
}

func InitNotificationServiceClientV2(c *config.Config) fcm.Client {
	client, err := fcm.NewClient(c.FirebaseApiKey)
	if err != nil {
		panic("Failed to create FCM client, err: " + err.Error())
	}

	return *client
}

type SendNotificationRequestV2 struct {
	To    string
	Title string
	Body  string
	Data  map[string]interface{}
}

func (s *NotificationServiceClientV2) SendNotificationV2(ctx context.Context, req *SendNotificationRequestV2) (*fcm.Response, error) {
	return s.Client.Send(&fcm.Message{
		To:       req.To,
		Priority: "high",
		Notification: &fcm.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Data: req.Data,
	})
}
