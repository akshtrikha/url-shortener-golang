//? This is the file which defines the route to hit to shorten the url

package routes

import (
	"time"

	"github.com/akshtrikha/url-shortener-golang/api/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/asaskevich/govalidator"
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

	return nil
}
