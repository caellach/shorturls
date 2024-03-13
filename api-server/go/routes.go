package main

// Loads all the routes for the application

import (
	"github.com/caellach/shorturl/api-server/go/controllers/auth"
	"github.com/caellach/shorturl/api-server/go/controllers/health"
	"github.com/caellach/shorturl/api-server/go/controllers/url"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoadRoutes loads all the routes for the application
func LoadRoutes(App *fiber.App, MongoClient *mongo.Client) {
	// Load the routes for the application
	_logger.Println("Loading routes...")

	// Health routes
	health.CreateHealthRoutes(App)

	// Auth routes
	auth.CreateAuthRoutes(App, MongoClient)

	// Url routes
	url.CreateUrlRoutes(App, MongoClient)

	_logger.Println("Routes loaded!")
}
