package firebase

import (
	"context"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"worker/config"
)

// Firebase is a struct that holds Firebase app instance
// and has function to easily access GCF features such as Firestore, Messaging, etc.
type Firebase struct {
	app *firebase.App
}

// FirestoreClient returns Firestore client instance
// client must be closed when we're finish using it with client.Close()
func (f *Firebase) FirestoreClient(ctx context.Context) (*firestore.Client, error) {
	client, err := f.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// MessagingClient returns FCM client instance
func (f *Firebase) MessagingClient(ctx context.Context) (*messaging.Client, error) {
	client, err := f.app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// PubSubClient return new instance of Pub/Sub client
func (f *Firebase) PubSubClient(ctx context.Context) (*pubsub.Client, error) {
	client, err := pubsub.NewClient(ctx, config.GCPProject)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// New create new Firebase
func New(cxt context.Context) (*Firebase, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Firebase{app}, nil
}
