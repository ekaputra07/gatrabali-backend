package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Delete all entries older than 3 month
func deleteOldEntries(ctx context.Context, client *firestore.Client) {
	threeMonthAgoMs := time.Now().AddDate(0, -3, 0).Unix() * 1000
	fmt.Printf("Deleting all entries older than %v\n", threeMonthAgoMs)
	iter := client.Collection("entries").
		Where("published_at", "<", threeMonthAgoMs).
		Limit(500).
		Documents(ctx)

	batch := client.Batch()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		batch.Delete(doc.Ref)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Migrate categories to category
func migrateCategory(ctx context.Context, client *firestore.Client) {
	iter := client.Collection("entries").
		OrderBy("published_at", firestore.Asc).
		// StartAfter(1561981998000).
		Limit(500).
		Documents(ctx)

	batch := client.Batch()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		data := doc.Data()
		fmt.Println(data["published_at"])

		pubAt := data["published_at"]
		categories := data["categories"].([]interface{})

		updatedFields := make(map[string]interface{})

		// This is to fix wrong data format
		if _, ok := pubAt.(int64); !ok {
			updatedFields["published_at"] = int64(pubAt.(float64))
		}

		if _, ok := categories[0].(int64); !ok {
			updatedFields["category_id"] = int64(categories[0].(float64))
		} else {
			updatedFields["category_id"] = categories[0].(int64)
		}
		batch.Set(doc.Ref, updatedFields, firestore.MergeAll)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "gatrabali")
	if err != nil {
		panic(err.Error())
	}

	// delete old entries
	// deleteOldEntries(ctx, client)
	migrateCategory(ctx, client)
}
