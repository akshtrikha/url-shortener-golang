// This file will be used to get the shortened url and return the original url to redirec to

package routes

import (
	"github.com/akshtrikha/url-shortener-golang/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL takes the shortened URL and redirects to the original URL
func ResolveURL(c *fiber.Ctx) error {
	// Get the url from the params in the fiber.Context
	url := c.Params("url")

	// Connect to the redis database using CreateClient implemented in the database package
	r := database.CreateClient(1)
	// Defer the close call to the redis client
	defer r.Close()

	// Get the value, err from the redis database
	// By passing the context to the database and the shortened URL
	value, err := r.Get(database.Ctx, url).Result()

	// Check for any error from the redis or any error occured while connecting to the redis
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shortened URL not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to the database",
		})
	}

	// Create another redis client which handle the counter
	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
