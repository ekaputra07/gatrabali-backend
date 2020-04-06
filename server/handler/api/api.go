package api

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber"

	"server/common/constant"
	"server/firebase"
)

// Handler represents the handler for APIs
type Handler struct {
	fb *firebase.Firebase
}

// New returns Handler instance
func New(fb *firebase.Firebase) *Handler {
	return &Handler{fb: fb}
}

func (h *Handler) firestore(ctx context.Context) (f *firestore.Client, err error) {
	f, err = h.fb.FirestoreClient(ctx)
	return
}

// Routes is collection handler for API
func (h *Handler) Routes(app *fiber.Fiber, pathPrefix string) {

	api := app.Group(pathPrefix)
	api.Get("/feeds", h.handleFeeds())

	api.Get("/entries", h.handleEntries(constant.Entries))
	api.Get("/entries/:entryId", h.handleEntry(constant.Entries))

	api.Get("/kriminal/entries", h.handleEntries(constant.Kriminal))
	api.Get("/kriminal/entries/:entryId", h.handleEntry(constant.Kriminal))

	api.Get("/baliunited/entries", h.handleEntries(constant.BaliUnited))
	api.Get("/baliunited/entries/:entryId", h.handleEntry(constant.BaliUnited))

	api.Get("/balebengong/entries", h.handleEntries(constant.BaleBengong))
	api.Get("/balebengong/entries/:entryId", h.handleEntry(constant.BaleBengong))

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
