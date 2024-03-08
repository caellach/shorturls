package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
)

// Make the mongo client a singleton
var client *mongo.Client

func GetMongoClient(MongoDBConfig *config.MongoDBConfig) *mongo.Client {
	if client != nil {
		return client
	}

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
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
