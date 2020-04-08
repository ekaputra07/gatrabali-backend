package sync

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"

	"server/common/constant"
	"server/common/types"
)

// storeCategories calls Miniflux categories API and store the objects to Firestore
func storeCategories(ctx context.Context, store *firestore.Client, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		categories, err := getCategories()
		if err != nil {
			return fmt.Errorf("storeCategories failed: %s", err)
		}

		// write in batch
		batch := store.Batch()
		for _, cat := range *categories {
			docRef := store.Collection(constant.Categories).Doc(strconv.FormatInt(cat.ID, 10))
			batch.Set(docRef, cat)
		}
		_, err = batch.Commit(ctx)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := store.Collection(constant.Categories).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeCategories: %v", *payload.Op)
}

// storeFeed calls Miniflux feeds API and store the object to Firestore
func storeFeed(ctx context.Context, store *firestore.Client, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		feed, err := getFeed(*payload.ID)
		if err != nil {
			return fmt.Errorf("storeFeed failed: %s", err)
		}
		_, err = store.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, feed)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := store.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeFeed: %v", *payload.Op)
}

// storeEntry calls Miniflux entries API and store the object to Firestore
func storeEntry(ctx context.Context, store *firestore.Client, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		entry, err := getEntry(*payload.ID)
		if err != nil {
			return fmt.Errorf("storeEntry failed: %s", err)
		}
		// if category `kriminal` or `baliunited` store so sparate collection
		if entry.CategoryID == 11 {
			entry.ID = entry.PublishedAt
			_, err = store.Collection(constant.Kriminal).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID == 12 {
			entry.ID = entry.PublishedAt
			_, err = store.Collection(constant.BaliUnited).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID > 12 {
			_, err = store.Collection(constant.BaleBengong).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		} else {
			_, err = store.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		}
		return err

	} else if *payload.Op == constant.OpDelete {
		// we don't support delete on separate collection for now eg. kriminal
		_, err := store.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeEntry: %v", *payload.Op)
}
