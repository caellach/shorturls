package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
)

var app *fiber.App
var discordConfig *config.DiscordConfig
var mongoClient *mongo.Client
var validAuthProviders = make(map[string]string)

var authStatesCollection *mongo.Collection
var usersCollection *mongo.Collection

var signingSecret []byte

// var logger = log.New(os.Stdout, "auth: ", log.LstdFlags)

func initProviders() {
	validAuthProviders["discord"] = "https://discord.com/oauth2/authorize?client_id=" + discordConfig.ClientID + "&response_type=code&scope=email+identify"

	//validAuthProviders["google"] = "https://accounts.google.com/o/oauth2/v2/auth"
	//validAuthProviders["facebook"] = "https://www.facebook.com/v12.0/dialog/oauth"
}

func CreateAuthRoutes(App *fiber.App, MongoClient *mongo.Client) {
	app = App
	discordConfig = &config.ServerConfig.Providers.DiscordConfig
	mongoClient = MongoClient

	authStatesCollection = mongoClient.Database("shared").Collection("authStates")
	usersCollection = mongoClient.Database("shared").Collection("users")

	initProviders()

	signingSecret = []byte(config.ServerConfig.Token.Secret)

	// Load the routes for the application
	app.Get("/api/auth/callback", authProviderCallback)
	app.Get("/api/auth/:authProvider", getAuthProviderOAuthURL)
}
