package health

import "github.com/gofiber/fiber/v2"

var app *fiber.App

func CreateHealthRoutes(App *fiber.App) {
	app = App

	// Load the routes for the application
	app.Get("/api/health", getHealth)
}
