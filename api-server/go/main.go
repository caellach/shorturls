package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
	"github.com/caellach/shorturl/api-server/go/pkg/mongo"
)

func main() {
	// load config
	_config := config.LoadConfig(config.DefaultConfigParams())

	app := fiber.New(
		fiber.Config{
			Prefork:           _config.App.Prefork,
			JSONDecoder:       sonic.Unmarshal,
			JSONEncoder:       sonic.Marshal,
			AppName:           _config.App.Name,
			Concurrency:       _config.App.Concurrency,
			EnablePrintRoutes: _config.App.EnablePrintRoutes,
		},
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization, Origin, Accept",
	}))

	app.Use(helmet.New())

	// Set default response headers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Server-Programming-Language", "Golang v1.22.0")
		return c.Next()
	})

	app.Use(monitor.New(monitor.Config{
		Title: _config.App.Name + " Metrics Page",
		Next: func(c *fiber.Ctx) bool {
			// only /metrics should show the metrics page
			return c.Path() != "/metrics"
		},
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	var mongoDBConfig = &config.MongoDBConfig{
		Uri:      _config.MongoDB.Uri,
		Port:     _config.MongoDB.Port,
		Database: _config.MongoDB.Database,
		Username: _config.MongoDB.Username,
		Password: _config.MongoDB.Password,
	}

	var mongoClient = mongo.GetMongoClient(mongoDBConfig)
	LoadRoutes(app, mongoClient)

	log.Fatal(app.Listen(":3000"))
}
