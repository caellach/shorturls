package url

import "go.mongodb.org/mongo-driver/bson/primitive"

type PutUrlRequest struct {
	Url string `json:"url"`
}

type UrlDocument struct {
	Id       string `json:"id" bson:"id"`
	UserId   string `json:"userId" bson:"userId"`
	Url      string `json:"url" bson:"url"`
	ShortUrl string `json:"shortUrl" bson:"-"`

	// DateTime fields
	Created primitive.DateTime `json:"created" bson:"created"`
	Updated primitive.DateTime `json:"updated,omitempty" bson:"updated,omitempty"`
	Deleted primitive.DateTime `json:"deleted,omitempty" bson:"deleted,omitempty"`
}
