package main

import (
	"context"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"server/common/constant"
	"google.golang.org/api/iterator"
)

// DB is represents the datastore, in this case its Firebase
type DB struct {
	app *firebase.App
}

// MakeDB returns instance of DB
func MakeDB(app *firebase.App) *DB {
	return &DB{app}
}

// GetFeeds returns all feeds
func (db *DB) GetFeeds(ctx context.Context) (*[]map[string]interface{}, error) {
	client, err := db.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	iter := client.Collection(constant.Feeds).Documents(ctx)
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
		feeds = append(feeds, map[string]interface{}{"id": data["id"], "title": data["title"]})
	}
	return &feeds, nil
}

// GetAllCategories returns all categories on database
func (db *DB) GetAllCategories(ctx context.Context) (*[]map[string]interface{}, error) {
	client, err := db.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	iter := client.Collection(constant.Categories).OrderBy("title", firestore.Asc).Documents(ctx)
	categories := []map[string]interface{}{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}
		categories = append(categories, doc.Data())
	}
	return &categories, nil
}

// GetCollectionEntry returns single entry from a specific collection
func (db *DB) GetCollectionEntry(ctx context.Context, collection string, id int) (*map[string]interface{}, error) {
	client, err := db.app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	doc, err := client.Collection(collection).Doc(strconv.Itoa(id)).Get(ctx)
	if err != nil {
		return nil, err
	}
	data := doc.Data()
	return &data, nil
}

// GetCollectionEntries returns paginated entries on specified collection
func (db *DB) GetCollectionEntries(
	ctx context.Context,
	collection, orderBy string,
	cursor, limit int) (*[]map[string]interface{}, error) {

	client, err := db.app.Firestore(ctx)
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

	query := client.Collection(collection).
		OrderBy(orderBy, firestore.Desc).
		Limit(limit)

	if cursor > 0 {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(ctx)
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

// GetCollectionCategoryEntries returns paginated entries in a category
func (db *DB) GetCollectionCategoryEntries(
	ctx context.Context,
	collection string,
	category int,
	orderBy string,
	cursor, limit int) (*[]map[string]interface{}, error) {

	client, err := db.app.Firestore(ctx)
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

	query := client.Collection(collection).
		Where("category_id", "==", category).
		OrderBy(orderBy, firestore.Desc).
		Limit(limit)

	if cursor > 0 {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(ctx)
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
