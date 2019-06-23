package function

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Firebase is a struct that holds Firebase app instance
// and has function to easily access GCF features such as Firestore, Messaging, etc.
type Firebase struct {
	app *firebase.App
}

var firebaseApp *Firebase

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

func init() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		panic("Unable to instantiate Firebase App")
	}
	firebaseApp = &Firebase{app}
}
