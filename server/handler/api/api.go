package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber"

	"server/common/constant"
	"server/firebase"
)

// Group is collection handler for API
func Group(
	ctx context.Context,
	pathPrefix string,
	app *fiber.Fiber,
	fb *firebase.Firebase) {

	api := app.Group(pathPrefix)
	api.Get("/feeds", handleFeeds(ctx, fb))

	api.Get("/entries", handleEntries(ctx, fb, constant.Entries))
	api.Get("/entries/:entryId", handleEntry(ctx, fb, constant.Entries))

	api.Get("/kriminal/entries", handleEntries(ctx, fb, constant.Kriminal))
	api.Get("/kriminal/entries/:entryId", handleEntry(ctx, fb, constant.Kriminal))

	api.Get("/baliunited/entries", handleEntries(ctx, fb, constant.BaliUnited))
	api.Get("/baliunited/entries/:entryId", handleEntry(ctx, fb, constant.BaliUnited))

	api.Get("/balebengong/entries", handleEntries(ctx, fb, constant.BaleBengong))
	api.Get("/balebengong/entries/:entryId", handleEntry(ctx, fb, constant.BaleBengong))

	// server error handler
	api.Use(func(c *fiber.Ctx) {
		if c.Error() != nil {
			log.Println("[ERROR]", c.Error())
			c.SendStatus(http.StatusInternalServerError)
			return
		}
		c.Next()
	})
}
