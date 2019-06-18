package gatrabali

import (
	"cloud.google.com/go/firestore"
	"context"
	"os"
)

// Firestore returns Firestore client instance
func Firestore() (*firestore.Client, error) {
	client, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT")) // GCP_PROJECT will avilable from Cloud Function environment
	if err != nil {
		return nil, err
	}
	return client, nil
}
