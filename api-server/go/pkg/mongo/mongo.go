package mongo

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func CreateCollectionsAndIndexes(client *mongo.Client) {
	if client == nil {
		_logger.Println("Mongo client not initialized")
		return
	}

	// Create the index for the url collection
	_logger.Println("Creating collections and indexes...")

	newMongoCollections := []NewMongoCollection{
		{
			Database:   "shared",
			Collection: "authStates",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.M{"state": 1},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys:    bson.M{"insertedAt": 1},
					Options: options.Index().SetExpireAfterSeconds(600),
				},
			},
		},
		{
			Database:   "shared",
			Collection: "users",
			Indexes: []mongo.IndexModel{
				{
					Keys:    map[string]interface{}{"email": 1},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{
						{Key: "providers[].id", Value: 1},
						{Key: "providers[].name", Value: 1},
					},
				},
			},
		},
		{
			Database:   "shorturls",
			Collection: "metadata",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.M{"userId": 1},
					Options: options.Index().SetUnique(true),
				},
			},
		},
		{
			Database:   "shorturls",
			Collection: "ogpData",
			Indexes:    []mongo.IndexModel{},
		},
		{
			Database:   "shorturls",
			Collection: "urls",
			Indexes: []mongo.IndexModel{
				{
					Keys:    bson.M{"id": 1},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.M{"userId": 1},
				},
			},
		},
	}

	if createCollectionsAndIndexes(client, newMongoCollections) {
		_logger.Println("Collections and indexes created!")
	} else {
		panic("Failed to create all collections and indexes")
	}
}

func createCollectionsAndIndexes(client *mongo.Client, collections []NewMongoCollection) bool {
	success := true
	for _, collection := range collections {
		// check if database exists
		_, err := client.ListDatabaseNames(context.Background(), bson.D{})
		if err != nil {
			_logger.Println("Failed to list databases")
			log.Fatal(err)
			success = false
		}

		database := client.Database(collection.Database)
		// check if collection exists
		collections, err := database.ListCollectionNames(context.Background(), bson.D{})
		if err != nil {
			_logger.Println("Failed to list collections")
			log.Fatal(err)
			success = false
		}

		// create collection if it doesn't exist
		if !contains(collections, collection.Collection) {
			err := database.CreateCollection(context.Background(), collection.Collection)
			if err != nil {
				_logger.Println("Failed to create collection for ", collection.Collection)
				log.Fatal(err)
				success = false
			}

			for _, index := range collection.Indexes {
				_, err := database.Collection(collection.Collection).Indexes().CreateOne(context.Background(), index)
				if err != nil {
					_logger.Println("Failed to create index", index, " for ", collection.Collection)
					log.Fatal(err)
					success = false
				}
			}

			_logger.Println("Indexes created!")
		}
	}

	return success
}
