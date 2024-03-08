package url

import (
	"github.com/caellach/shorturl/api-server/go/pkg/config"
	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/caellach/shorturl/api-server/go/pkg/wordlist"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var app *fiber.App
var mongoClient *mongo.Client
var wordList *wordlist.WordList

var shorturlsCollection *mongo.Collection

func CreateUrlRoutes(WordListConfig *config.WordListConfig, App *fiber.App, MongoClient *mongo.Client) {
	app = App
	mongoClient = MongoClient

	shorturlsCollection = mongoClient.Database("shorturls").Collection("urls")

	// Load the word list
	wordList = wordlist.LoadWordList(WordListConfig)

	// Load the routes for the application
	// Public routes
	app.Get("/u/:id", getUrlById)

	// Authenticated routes
	app.Get("/u/", middleware.AuthRequired(), getUrls)
	app.Put("/u/", middleware.AuthRequired(), putUrl)
	app.Delete("/u/:id", middleware.AuthRequired(), deleteUrl)

}
