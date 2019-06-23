package function

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Firestore returns Firestore client instance
func Firestore() (*firestore.Client, error) {
	client, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT")) // GCP_PROJECT will avilable from Cloud Function environment
	if err != nil {
		return nil, err
	}
	return client, nil
}

// MessagingClient returns instance of FCM messaging client
func MessagingClient(ctx context.Context) (*messaging.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
