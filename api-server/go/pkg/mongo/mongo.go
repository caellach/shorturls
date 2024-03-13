package mongo

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
)

var _logger = log.New(os.Stdout, "mongo: ", log.LstdFlags)

// Make the mongo client a singleton
var client *mongo.Client

func GetMongoClient(MongoDBConfig *config.MongoDBConfig) *mongo.Client {
	if client != nil {
		return client
	}

	_logger.Println("Connecting to MongoDB...")

	// Set client options
	clientOptions := options.Client().ApplyURI(MongoDBConfig.Uri)
	clientOptions.SetAuth(options.Credential{
		Username:      MongoDBConfig.Username,
		Password:      MongoDBConfig.Password,
		AuthMechanism: "SCRAM-SHA-1",
	})
	clientOptions.ConnectTimeout = &[]time.Duration{5 * time.Second}[0]

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		_logger.Println("Failed to connect to MongoDB")
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		_logger.Println("Failed to ping MongoDB")
		log.Fatal(err)
	}

	_logger.Println("Connected to MongoDB!")

	return client
}
