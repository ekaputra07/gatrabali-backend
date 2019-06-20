package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"

	"github.com/apps4bali/gatrabali-backend/common/constant"
)

// DB is represents the datastore, in this case its Firebase
type DB struct {
	app *firebase.App
}

// MakeDB returns instance of DB
func MakeDB(app *firebase.App) *DB {
	return &DB{app}
}

// GetFeeds returns list feeds
func (db *DB) GetFeeds() (*[]map[string]interface{}, error) {
	ctx := context.Background()
	client, err := db.app.Firestore(ctx)
	if err != nil {
		return nil, err
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
	return &feeds, nil
}
