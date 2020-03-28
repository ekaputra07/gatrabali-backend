package firebase

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Firebase is a struct that holds Firebase app instance
// and has function to easily access GCF features such as Firestore, Messaging, etc.
type Firebase struct {
	app *firebase.App

	firestoreClientOnce sync.Once
	messagingClientOnce sync.Once
	pubsubClientOnce    sync.Once
	firestoreClient     *firestore.Client
	messagingClient     *messaging.Client
	pubsubClient        *pubsub.Client
}

// FirestoreClient returns Firestore client instance (lazily loaded)
func (f *Firebase) FirestoreClient(ctx context.Context) (*firestore.Client, error) {
	var err error
	f.firestoreClientOnce.Do(func() {
		f.firestoreClient, err = f.app.Firestore(ctx)
	})
	return f.firestoreClient, err
}

// MessagingClient returns FCM client instance (lazily loaded)
func (f *Firebase) MessagingClient(ctx context.Context) (*messaging.Client, error) {
	var err error
	f.messagingClientOnce.Do(func() {
		f.messagingClient, err = f.app.Messaging(ctx)
	})
	return f.messagingClient, err
}

// PubSubClient return new instance of Pub/Sub client (lazily loaded)
func (f *Firebase) PubSubClient(ctx context.Context, gcpProject string) (*pubsub.Client, error) {
	var err error
	f.pubsubClientOnce.Do(func() {
		f.pubsubClient, err = pubsub.NewClient(ctx, gcpProject)
	})
	return f.pubsubClient, err
}

// New create new Firebase
func New(cxt context.Context) (*Firebase, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Firebase{app: app}, nil
}
