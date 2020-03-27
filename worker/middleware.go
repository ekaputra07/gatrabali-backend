package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber"
)

// softErrorHandler is a middleware to handle error but NOT actually return error response.
// This is to avoid PubSub retrying.
func softErrorHandler() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		if c.Error() != nil {
			log.Println("[ERROR]", c.Error())
			c.SendStatus(http.StatusOK)
		}
		c.Next()
	}
}
