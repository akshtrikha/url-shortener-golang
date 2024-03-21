//? This is the file which defines the route to hit to shorten the url

package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/akshtrikha/url-shortener-golang/database"
	"github.com/akshtrikha/url-shortener-golang/logger"

	"github.com/akshtrikha/url-shortener-golang/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	Expiry          time.Duration `json:"expiry"`           // Hours remaining before the Short URL expires
	XRateRemaining  int           `json:"rate_limit"`       // Further number of request allowed before Rate Limiter reset
	XRateLimitReset time.Duration `json:"rate_limit_reset"` // Seconds remaining for the rate limiter to reset
}

// ShortenURL function parses the payload and returns the shortened URL
func ShortenURL(ctx *fiber.Ctx) error {
	logger.Log.Println("Shorten URL started")

	body := new(request)

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse the JSON payload"})
	}

	//? implement rate limiting
	// r = redisClient for rate limiting
	r := database.CreateClient(1)
	logger.Log.Println("Created the redis client")
	defer r.Close()
	value, err := r.Get(database.Ctx, ctx.IP()).Result()
	if err != nil {
		logger.Log.Println("redis error: ", err)
	}
	logger.Log.Println("value: ", value)
	if err == redis.Nil {
		_ = r.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		logger.Log.Println("IP already present in the database");
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			logger.Log.Println("ValInt is less than equal to 0 ", valInt)
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

	//? Custom shortened URL
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// r2 = redisClient for handling short and original urls
	r2 := database.CreateClient(0)
	defer r2.Close()
	val, _ := r2.Get(database.Ctx, id).Result()

	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom URL is already in use",
		})
	}

	//? Set the expiration of the URL to 24 hours
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	logger.Log.Println("body.URL: ", body.URL, "id: ", id)
	err = r2.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to the server",
		})
	}

	//? Creating the response object
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	//? Decrement the value of the allowed request to handle rate limiting
	r.Decr(database.Ctx, ctx.IP())

	val, _ = r.Get(database.Ctx, ctx.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r.TTL(database.Ctx, ctx.IP()).Result()
	resp.XRateLimitReset = ttl

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
