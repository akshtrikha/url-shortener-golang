//? This is the file which defines the route to hit to shorten the url

package routes

import (
	"github.com/akshtrikha/url-shortener-golang/database"
	"os"
	"strconv"
	"time"

	"github.com/akshtrikha/url-shortener-golang/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ? Defines the structure of a request
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// ? Defines the structure of a response
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL function parses the payload and returns the shortened URL
func ShortenURL(ctx *fiber.Ctx) error {

	body := new(request)

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse the JSON payload"})
	}

	//TODO: implement rate limiting
	r := database.CreateClient(1)
	defer r.Close()
	value, err := r.Get(database.Ctx, ctx.IP()).Result()
	if err == redis.Nil {
		_ = r.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			limit, _ := r.TTL(database.Ctx, ctx.IP()).Result()
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit,
			})
		}
	}

	//? validate URL
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	//? check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Invalid domain"})
	}

	//? enforce HTTPS
	body.URL = helpers.EnforceHTTP(body.URL)

	//? Decrement the value of the allowed request to handle rate limiting
	r.Decr(database.Ctx, ctx.IP())

	return nil
}
