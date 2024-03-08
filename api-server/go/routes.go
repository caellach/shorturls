package main

// Loads all the routes for the application

import (
	"github.com/caellach/shorturl/api-server/go/controllers/auth"
	"github.com/caellach/shorturl/api-server/go/controllers/health"
	"github.com/caellach/shorturl/api-server/go/controllers/url"
	"github.com/caellach/shorturl/api-server/go/pkg/config"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoadRoutes loads all the routes for the application
func LoadRoutes(Config *config.Config, App *fiber.App, MongoClient *mongo.Client) {
	// Load the routes for the application

	// Health routes
	health.CreateHealthRoutes(App)

	// Auth routes
	auth.CreateAuthRoutes(&Config.Providers.DiscordConfig, App, MongoClient)

	// Url routes
	url.CreateUrlRoutes(&Config.WordList, App, MongoClient)
}
