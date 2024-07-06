package firebase

import (
	"context"
	fcm "github.com/appleboy/go-fcm"
	"os"
)

var FCMClient *fcm.Client

func InitFirebase() error {
	var err error
	FCMClient, err = fcm.NewClient(
		context.Background(),
		fcm.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS")),
	)
	return err
}
