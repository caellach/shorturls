package url

import "go.mongodb.org/mongo-driver/bson/primitive"

type PutUrlRequest struct {
	Url string `json:"url"`
}

type UrlDocument struct {
	Id       string             `json:"id" bson:"id"`
	UserId   string             `json:"user_id" bson:"user_id"`
	Url      string             `json:"url" bson:"url"`
	ShortUrl string             `json:"shortUrl" bson:"-"`
	UseCount int                `json:"use_count" bson:"use_count"`
	LastUsed primitive.DateTime `json:"last_used,omitempty" bson:"last_used,omitempty"`

	// DateTime fields
	Created primitive.DateTime `json:"created" bson:"created"`
	Updated primitive.DateTime `json:"updated,omitempty" bson:"updated,omitempty"`
	Deleted primitive.DateTime `json:"deleted,omitempty" bson:"deleted,omitempty"`
}
