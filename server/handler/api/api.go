package api

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber"

	"server/common/constant"
	"server/common/service"
)

// Handler represents the handler for APIs
type Handler struct {
	google *service.Google
}

// New returns Handler instance
func New(google *service.Google) *Handler {
	return &Handler{google}
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
