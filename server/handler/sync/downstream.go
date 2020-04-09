package sync

import (
	"context"
	"fmt"
	"strconv"

	"server/common/constant"
	"server/common/types"
)

// storeCategories calls Miniflux categories API and store the objects to Firestore
func (h *Handler) storeCategories(ctx context.Context, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		categories, err := getCategories(ctx)
		if err != nil {
			return fmt.Errorf("storeCategories failed: %s", err)
		}

		// write in batch
		batch := h.google.Firestore.Batch()
		for _, cat := range *categories {
			docRef := h.google.Firestore.Collection(constant.Categories).Doc(strconv.FormatInt(cat.ID, 10))
			batch.Set(docRef, cat)
		}
		_, err = batch.Commit(ctx)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := h.google.Firestore.Collection(constant.Categories).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeCategories: %v", *payload.Op)
}

// storeFeed calls Miniflux feeds API and store the object to Firestore
func (h *Handler) storeFeed(ctx context.Context, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		feed, err := getFeed(ctx, *payload.ID)
		if err != nil {
			return fmt.Errorf("storeFeed failed: %s", err)
		}
		_, err = h.google.Firestore.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, feed)
		return err

	} else if *payload.Op == constant.OpDelete {
		_, err := h.google.Firestore.Collection(constant.Feeds).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeFeed: %v", *payload.Op)
}

// storeEntry calls Miniflux entries API and store the object to Firestore
func (h *Handler) storeEntry(ctx context.Context, payload *types.SyncPayload) error {
	if *payload.Op == constant.OpWrite {
		entry, err := getEntry(ctx, *payload.ID)
		if err != nil {
			return fmt.Errorf("storeEntry failed: %s", err)
		}
		// if category `kriminal` or `baliunited` store so sparate collection
		if entry.CategoryID == 11 {
			entry.ID = entry.PublishedAt
			_, err = h.google.Firestore.Collection(constant.Kriminal).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID == 12 {
			entry.ID = entry.PublishedAt
			_, err = h.google.Firestore.Collection(constant.BaliUnited).Doc(strconv.FormatInt(entry.ID, 10)).Set(ctx, entry)
		} else if entry.CategoryID > 12 {
			_, err = h.google.Firestore.Collection(constant.BaleBengong).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		} else {
			_, err = h.google.Firestore.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Set(ctx, entry)
		}
		return err

	} else if *payload.Op == constant.OpDelete {
		// we don't support delete on separate collection for now eg. kriminal
		_, err := h.google.Firestore.Collection(constant.Entries).Doc(strconv.FormatInt(*payload.ID, 10)).Delete(ctx)
		return err
	}
	return fmt.Errorf("Invalid operation for storeEntry: %v", *payload.Op)
}
