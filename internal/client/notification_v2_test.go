package client

import (
	"context"
	"os"
	"testing"

	"github.com/lntvan166/e2tech-booking-svc/internal/config"
	"github.com/stretchr/testify/require"
)

const (
	FIREBASE_API_KEY_TEST = "AAAAcgxFAtE:APA91bGvqoDv5GTYpZqqND-VCryhzuwhIMHidvmnCgNEIUQhKV7a8M6oBgkqEHZib918KSFMwPT3c9i6_pXb50NE5kG8HUvEJnJXo4Qo6PkA3mIlEkOBvlJ2z2qFE5wozRD2A8tP9nIx"
	DEVICE_TOKEN_TEST     = "cE5iRQt8ReC3KL3fLdAgP4:APA91bGKSI7btxighoYeRUtZzh3b4_ns4E6HSrqhiThmjR1a6qP_hGysN15_VVbo_de6SkH21LhTyJNxP20O97IkcIqnrsHryXvV46Yh-2eBLJ7EeFHKkCFIPYr4b3xxlhy4DvAI2I-n"
)

var NotificationSvc *NotificationServiceClientV2

func TestMain(m *testing.M) {
	c := config.Config{
		FirebaseApiKey: FIREBASE_API_KEY_TEST,
	}

	NotificationSvc = &NotificationServiceClientV2{
		Client: InitNotificationServiceClientV2(&c),
	}

	os.Exit(m.Run())

}

func TestSendNotificationV2(t *testing.T) {
	r, err := NotificationSvc.SendNotificationV2(context.Background(), &SendNotificationRequestV2{
		To:    "cE5iRQt8ReC3KL3fLdAgP4:APA91bGKSI7btxighoYeRUtZzh3b4_ns4E6HSrqhiThmjR1a6qP_hGysN15_VVbo_de6SkH21LhTyJNxP20O97IkcIqnrsHryXvV46Yh-2eBLJ7EeFHKkCFIPYr4b3xxlhy4DvAI2I-n",
		Title: "Hello",
		Body:  "Hello World",
		Data:  map[string]interface{}{"phone": "0892892832"},
	})
	require.NoError(t, err)
	require.Equal(t, r.Success, 1)
}
