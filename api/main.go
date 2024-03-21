package main

import (
	"fmt"
	// "log"
	"os"

	"github.com/akshtrikha/url-shortener-golang/routes"
	"github.com/akshtrikha/url-shortener-golang/logger"
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	logger.Log.Println("setting up the routes")
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	logger.Log.Println("Main()")
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading environment variables:", err)
	}

	app := fiber.New()
	logger.Log.Println("fiber app created")

	// app.Use(logger.New())

	setupRoutes(app)
	logger.Log.Println("app is setup")

	// log.Fatal(app.Listen("localhost:" + os.Getenv("APP_PORT")))
	app.Listen("localhost:" + os.Getenv("APP_PORT"))
}
