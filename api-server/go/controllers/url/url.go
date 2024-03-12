package url

import (
	"github.com/caellach/shorturl/api-server/go/pkg/config"
	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/caellach/shorturl/api-server/go/pkg/wordlist"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var app *fiber.App
var mongoClient *mongo.Client
var wordList *wordlist.WordList

var metadataCollection *mongo.Collection
var shorturlsCollection *mongo.Collection

var websocketConnections = make(map[string][]*websocket.Conn)

func CreateUrlRoutes(App *fiber.App, MongoClient *mongo.Client) {
	app = App
	mongoClient = MongoClient

	metadataCollection = mongoClient.Database("shorturls").Collection("metadata")
	shorturlsCollection = mongoClient.Database("shorturls").Collection("urls")

	// Load the word list
	wordList = wordlist.LoadWordList(&config.ServerConfig.WordList)

	// Authenticated routes
	app.Get("/u/metadata", middleware.AuthRequired(), getUserMetadata)
	app.Get("/u/", middleware.AuthRequired(), getUrls)
	app.Put("/u/", middleware.AuthRequired(), putUrl)
	app.Delete("/u/:id", middleware.AuthRequired(), deleteUrl)

	// websocket route
	app.Use("/u/ws", websocket.New(urlWs))

	// Load the routes for the application
	// Public routes
	app.Get("/u/:id", redirectUrlById)
}
