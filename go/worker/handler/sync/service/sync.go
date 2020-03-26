package service

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/apps4bali/gatrabali-backend/go/common"
)

// StartCategorySync calls Miniflux categories API and store the objects to Firestore
func StartCategorySync(ctx context.Context, store *firestore.Client, payload *common.SyncPayload) error {
	if *payload.Op == common.OpWrite {
		categories, err := GetCategories()
		if err != nil {
			return fmt.Errorf("StartCategorySync failed: %s", err)
		}

		// write in batch
		batch := store.Batch()
		for _, cat := range *categories {
			docRef := store.Collection(common.Categories).Doc(strconv.FormatInt(cat.ID, 10))
			batch.Set(docRef, cat)
		}
		_, err = batch.Commit(ctx)
		return err

	} else if *payload.Op == common.OpDelete {
		_, err := store.Collection(common.Categories).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartCategorySync: %v", *payload.Op)
}

// StartFeedSync calls Miniflux feeds API and store the object to Firestore
func StartFeedSync(ctx context.Context, store *firestore.Client, payload *common.SyncPayload) error {
	if *payload.Op == common.OpWrite {
		feed, err := GetFeed(*payload.ID)
		if err != nil {
			return fmt.Errorf("StartFeedSync failed: %s", err)
		}
		_, err = store.Collection(common.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, feed)
		return err

	} else if *payload.Op == common.OpDelete {
		_, err := store.Collection(common.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartFeedSync: %v", *payload.Op)
}

// StartEntrySync calls Miniflux entries API and store the object to Firestore
func StartEntrySync(ctx context.Context, store *firestore.Client, payload *common.SyncPayload) error {
	if *payload.Op == common.OpWrite {
		entry, err := GetEntry(*payload.ID)
		if err != nil {
			return fmt.Errorf("StartEntrySync failed: %s", err)
		}
		// if category `kriminal` or `baliunited` store so sparate collection
		if entry.CategoryID == 11 {
			entry.ID = entry.PublishedAt
			_, err = store.Collection(common.Kriminal).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID == 12 {
			entry.ID = entry.PublishedAt
			_, err = store.Collection(common.BaliUnited).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID > 12 {
			_, err = store.Collection(common.BaleBengong).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		} else {
			_, err = store.Collection(common.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		}
		return err

	} else if *payload.Op == common.OpDelete {
		// we don't support delete on separate collection for now eg. kriminal
		_, err := store.Collection(common.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartEntrySync: %v", *payload.Op)
}
