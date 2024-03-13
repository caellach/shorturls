package mongo

import "go.mongodb.org/mongo-driver/mongo"

type NewMongoCollection struct {
	Database   string
	Collection string
	Indexes    []mongo.IndexModel
}
