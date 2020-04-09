package main

import (
	"context"
	"log"

	"github.com/fiberweb/apikey"
	pubs "github.com/fiberweb/pubsub"
	"github.com/gofiber/fiber"

	"server/common/service"
	"server/config"
	"server/handler/api"
	"server/handler/events"
	"server/handler/push"
	"server/handler/sync"
)

var gcp *service.Google

func main() {
	ctx := context.Background()

	// initialize Firebase app
	var err error
	gcp, err = service.NewGoogle(ctx, config.GCPProject)
	if err != nil {
		log.Fatalln("Unable to initialize Firebase app:", err)
	}

	app := fiber.New()

	// all /pubsub/** are to handle PubSub requests (protected by api key)
	pubsub := app.Group("/pubsub")
	pubsub.Use(apikey.New(apikey.Config{
		Key: config.PubSubAPIKey,
		Skip: func(c *fiber.Ctx) bool {
			if "dev" == config.PubSubAPIKey {
				return true
			}
			return false
		},
	}))

	pubsub.Use(pubs.New(pubs.Config{Debug: false})) // pubsub middleware
	pubsub.Post("/sync-data", sync.New(gcp).Handle())
	pubsub.Post("/push-notification", push.New(gcp).Handle())
	pubsub.Post("/firestore-events", events.New(gcp).Handle())
	pubsub.Use(softErrorHandler()) // always return OK response to avoid PubSub retrying

	// all /api/** are to REST apis for clients
	api.New(gcp).Routes(app, "/api/v1")

	app.Listen(config.ServicePort)
}
