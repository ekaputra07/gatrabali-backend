package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber"

	"server/firebase"
)

const (
	defaultMaxAge  = 900  // browser cache: 15 mins
	defaultSmaxAge = 3600 // CDN cache: 1 hr
)

func setCacheControl(c *fiber.Ctx, maxAge, smaxAge int64) {
	c.Set("Cache-Control", fmt.Sprintf("public, max-age=%v, s-maxage=%v", maxAge, smaxAge))
}

// this is a workaround to Fiber's c.JSON() content-type doesn't include charset by default.
// https://github.com/gofiber/fiber/issues/248
func sendJSON(c *fiber.Ctx, body interface{}) {
	c.JSON(body)
	c.Set("Content-type", "application/json; charset=utf-8")
}

func handleFeeds(ctx context.Context, fb *firebase.Firebase) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		// get firestore
		firestore, err := fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}

		feeds := getFeeds(context.Background(), firestore)
		if len(feeds) > 0 {
			setCacheControl(c, defaultMaxAge, defaultSmaxAge)
		}
		sendJSON(c, feeds)
	}
}

func handleEntries(ctx context.Context, fb *firebase.Firebase, collection string) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		cat, err := strconv.Atoi(c.Query("categoryId"))
		if err != nil {
			cat = 0
		}

		lim, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			lim = 10
		}

		cur, err := strconv.Atoi(c.Query("cursor"))
		if err != nil {
			cur = 0
		}

		// get firestore
		firestore, err := fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}

		var entries []map[string]interface{}

		opts := queryopts{
			Collection: collection,
			Cursor:     cur,
			Limit:      lim,
		}
		if cat > 0 {
			opts.Category = cat
		}

		entries = getEntries(context.Background(), firestore, opts)
		if len(entries) > 0 {
			setCacheControl(c, defaultMaxAge, defaultSmaxAge)
		}
		sendJSON(c, entries)
	}
}

func handleEntry(ctx context.Context, fb *firebase.Firebase, collection string) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		id := c.Params("entryId")

		// get firestore
		firestore, err := fb.FirestoreClient(ctx)
		if err != nil {
			c.Next(err)
			return
		}
		entry, err := getEntry(context.Background(), firestore, queryopts{Collection: collection, ID: id})
		if err != nil {
			c.SendStatus(http.StatusNotFound)
			return
		}
		setCacheControl(c, defaultMaxAge, 86400) // 24 hr since individual entry won't change much
		sendJSON(c, entry)
	}
}
