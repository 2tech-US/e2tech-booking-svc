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
	DEVICE_TOKEN_TEST     = "dm38o087QRm6Bu8jTXB5Qq:APA91bGjz940Ynf-49xwX0wl0vn67QL3pNKpiKlUJ_06NV5KtOYxSMSY5RRMS3vVgIKU1x22sPcE63OLf5ZQK_ttkD4xt1Ved4yaNcSitxPN9GRf2DKWTD6pA8R3vz1o1ju8RAsyDee3"
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
		To:    "dm38o087QRm6Bu8jTXB5Qq:APA91bGjz940Ynf-49xwX0wl0vn67QL3pNKpiKlUJ_06NV5KtOYxSMSY5RRMS3vVgIKU1x22sPcE63OLf5ZQK_ttkD4xt1Ved4yaNcSitxPN9GRf2DKWTD6pA8R3vz1o1ju8RAsyDee3",
		Title: "Hello",
		Body:  "Hello World",
		Data:  map[string]interface{}{"phone": "0892892832"},
	})
	require.NoError(t, err)
	require.Equal(t, r.Success, 1)
}
