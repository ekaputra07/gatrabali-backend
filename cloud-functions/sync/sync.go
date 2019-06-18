package sync

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"strconv"

	"gatrabali/constant"
	"gatrabali/model"
)

// StartCategorySync calls Miniflux categories API and store the objects to Firestore
func StartCategorySync(store *firestore.Client, payload *model.Payload) error {
	ctx := context.Background()

	if *payload.Op == constant.OpWrite {
		categories, err := GetCategories()
		if err != nil {
			return err
		}

		// write in batch
		batch := store.Batch()
		for _, cat := range *categories {
			docRef := store.Collection(constant.Categories).Doc(strconv.FormatInt(cat.ID, 10))
			m, err := cat.ToMap()
			if err == nil {
				batch.Set(docRef, m)
			}
		}
		_, err = batch.Commit(ctx)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := store.Collection(constant.Categories).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartCategorySync: %v", *payload.Op)
}

// StartFeedSync calls Miniflux feeds API and store the object to Firestore
func StartFeedSync(store *firestore.Client, payload *model.Payload) error {
	ctx := context.Background()

	if *payload.Op == constant.OpWrite {
		feed, err := GetFeed(*payload.ID)
		if err != nil {
			return err
		}
		m, err := feed.ToMap()
		if err != nil {
			return err
		}
		_, err = store.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, m)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := store.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartFeedSync: %v", *payload.Op)
}

// StartEntrySync calls Miniflux entries API and store the object to Firestore
func StartEntrySync(store *firestore.Client, payload *model.Payload) error {
	ctx := context.Background()

	if *payload.Op == constant.OpWrite {
		entry, err := GetEntry(*payload.ID)
		if err != nil {
			return err
		}
		m, err := entry.ToMap()
		if err != nil {
			return err
		}
		_, err = store.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, m)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := store.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for StartEntrySync: %v", *payload.Op)
}
