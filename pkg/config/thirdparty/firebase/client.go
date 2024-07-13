package firebase

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/appleboy/go-fcm"
	"os"
)

var FCMClient *fcm.Client

func InitFirebase() error {
	encodedCredentials := os.Getenv("FIREBASE_CREDENTIALS")
	if encodedCredentials == "" {
		return errors.New("FIREBASE_CREDENTIALS environment variable is not set")
	}

	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp("", "firebase-credentials-*.json")
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			err = errors.New("Failed to remove temporary file")
		}
	}(tmpFile.Name())

	if _, err := tmpFile.Write(decodedCredentials); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	FCMClient, err = fcm.NewClient(
		context.Background(),
		fcm.WithCredentialsFile(tmpFile.Name()),
	)
	return err
}
