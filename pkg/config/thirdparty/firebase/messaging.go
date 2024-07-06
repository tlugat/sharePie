package firebase

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"log"
)

// SendNotification sends a notification to a list of tokens using Firebase Cloud Messaging.
func SendNotification(tokens []*string, notification messaging.Notification) error {
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens provided")
	}

	ctx := context.Background()
	var validTokens []string
	for _, token := range tokens {
		if token != nil && *token != "" {
			validTokens = append(validTokens, *token)
		}
	}

	if len(validTokens) == 0 {
		return fmt.Errorf("no valid tokens provided")
	}

	message := &messaging.MulticastMessage{
		Tokens:       validTokens,
		Notification: &notification,
	}

	response, err := FCMClient.SendMulticast(ctx, message)
	if err != nil {
		log.Printf("Failed to send notifications: %v", err)
		return err
	}

	for idx, res := range response.Responses {
		if res.Success {
			log.Printf("Successfully sent notification to token %s", validTokens[idx])
		} else {
			log.Printf("Failed to send notification to token %s: %v", validTokens[idx], res.Error)
		}
	}

	if response.FailureCount > 0 {
		return fmt.Errorf("encountered errors sending notifications: %d failures", response.FailureCount)
	}
	return nil
}
