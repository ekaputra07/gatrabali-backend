package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"

	"github.com/apps4bali/gatrabali-backend/common/constant"
)

// DB is represents the datastore, in this case its Firebase
type DB struct {
	app *firebase.App
	ctx context.Context
}

// MakeDB returns instance of DB
func MakeDB(app *firebase.App) *DB {
	return &DB{app, context.Background()}
}

// GetFeeds returns all feeds
func (db *DB) GetFeeds() (*[]map[string]interface{}, error) {
	client, err := db.app.Firestore(db.ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	iter := client.Collection(constant.Feeds).Documents(db.ctx)
	feeds := []map[string]interface{}{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		data := doc.Data()
		feeds = append(feeds, map[string]interface{}{"id": data["id"], "name": data["title"]})
	}
	return &feeds, nil
}

// GetAllEntries returns paginated entries
func (db *DB) GetAllEntries(cursor, limit int) (*[]map[string]interface{}, error) {
	client, err := db.app.Firestore(db.ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// default limit 10
	// max limit 20
	if limit == 0 {
		limit = 10
	} else if limit > 20 {
		limit = 20
	}

	query := client.Collection(constant.Entries).
		OrderBy("id", firestore.Desc).
		Limit(limit)

	if cursor > 0 {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(db.ctx)
	entries := []map[string]interface{}{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		entries = append(entries, doc.Data())
	}
	return &entries, nil
}

// GetCategoryEntries returns paginated entries in a category
func (db *DB) GetCategoryEntries(category, cursor, limit int) (*[]map[string]interface{}, error) {
	client, err := db.app.Firestore(db.ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// default limit 10
	// max limit 20
	if limit == 0 {
		limit = 10
	} else if limit > 20 {
		limit = 20
	}

	query := client.Collection(constant.Entries).
		Where("categories", "array-contains", category).
		OrderBy("id", firestore.Desc).
		Limit(limit)

	if cursor > 0 {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(db.ctx)
	entries := []map[string]interface{}{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		entries = append(entries, doc.Data())
	}
	return &entries, nil
}
