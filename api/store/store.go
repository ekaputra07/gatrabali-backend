package store

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/apps4bali/gatrabali-backend/common/constant"
	"google.golang.org/api/iterator"
	// "github.com/apps4bali/gatrabali-backend/common/model"
	"log"
)

// FirebaseApp return default Firebase app instance
var FirebaseApp *firebase.App
var ctx = context.Background()

func init() {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	FirebaseApp = app
}

// Firestore return Firestore client instance
// IMPORTANT: client should be closed when finished using it.
func Firestore() (*firestore.Client, error) {
	client, err := FirebaseApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, err
}

// GetFeeds return a list of Feed object
func GetFeeds() *[]map[string]interface{} {
	client, err := Firestore()
	if err != nil {
		log.Println("error initializing Firestore client: %v\n", err)
		return nil
	}
	defer client.Close()

	iter := client.Collection(constant.Feeds).Documents(ctx)
	var feeds []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		}
		data := doc.Data()
		feeds = append(feeds, map[string]interface{}{"id": data["id"], "name": data["name"]})
	}
	return &feeds
}
