package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber"
)

const (
	defaultMaxAge  = 900  // browser cache: 15 mins
	defaultSmaxAge = 3600 // CDN cache: 1 hr
)

func (h *Handler) setCacheControl(c *fiber.Ctx, maxAge, smaxAge int64) {
	c.Set("Cache-Control", fmt.Sprintf("public, max-age=%v, s-maxage=%v", maxAge, smaxAge))
}

// this is a workaround to Fiber's c.JSON() content-type doesn't include charset by default.
// https://github.com/gofiber/fiber/issues/248
func (h *Handler) sendJSON(c *fiber.Ctx, body interface{}) {
	c.JSON(body)
	c.Set("Content-type", "application/json; charset=utf-8")
}

func (h *Handler) handleFeeds() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		feeds, err := h.getFeeds(context.Background())
		if err != nil {
			c.Next(err)
			return
		}

		if len(feeds) > 0 {
			h.setCacheControl(c, defaultMaxAge, defaultSmaxAge)
		}
		h.sendJSON(c, feeds)
	}
}

func (h *Handler) handleEntries(collection string) func(*fiber.Ctx) {
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

		opts := queryopts{
			Collection: collection,
			Cursor:     cur,
			Limit:      lim,
		}
		if cat > 0 {
			opts.Category = cat
		}

		entries, err := h.getEntries(context.Background(), opts)
		if err != nil {
			c.Next(err)
			return
		}

		if len(entries) > 0 {
			h.setCacheControl(c, defaultMaxAge, defaultSmaxAge)
		}
		h.sendJSON(c, entries)
	}
}

func (h *Handler) handleEntry(collection string) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		id := c.Params("entryId")

		entry, err := h.getEntry(context.Background(), queryopts{Collection: collection, ID: id})
		if err != nil {
			c.SendStatus(http.StatusNotFound)
			return
		}
		h.setCacheControl(c, defaultMaxAge, 86400) // 24 hr since individual entry won't change much
		h.sendJSON(c, entry)
	}
}
