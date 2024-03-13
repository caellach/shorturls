package url

import "go.mongodb.org/mongo-driver/bson/primitive"

type OgpData struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	SiteName    string             `json:"siteName" bson:"siteName"`
	Title       string             `json:"title" bson:"title"`
	Type        string             `json:"type" bson:"type"`
	Url         string             `json:"url" bson:"url"`
}

type PutUrlRequest struct {
	Url string `json:"url"`
}

type UrlDocument struct {
	Id        string             `json:"id" bson:"id"`
	UserId    string             `json:"userId" bson:"userId"`
	Url       string             `json:"url" bson:"url"`
	ShortUrl  string             `json:"shortUrl" bson:"-"`
	UseCount  int                `json:"useCount" bson:"useCount"`
	LastUsed  primitive.DateTime `json:"lastUsed,omitempty" bson:"lastUsed,omitempty"`
	OgpDataId primitive.ObjectID `bson:"ogpDataId,omitempty"`

	// DateTime fields
	Created primitive.DateTime `json:"created" bson:"created"`
	Updated primitive.DateTime `json:"updated,omitempty" bson:"updated,omitempty"`
	Deleted primitive.DateTime `json:"deleted,omitempty" bson:"deleted,omitempty"`
}

type UserUrlMetadata struct {
	UserId       string             `json:"userId" bson:"userId"`
	ActiveCount  int                `json:"activeCount" bson:"activeCount"`
	CreatedCount int                `json:"createdCount" bson:"createdCount"`
	LastCreated  primitive.DateTime `json:"lastCreated" bson:"lastCreated"`
}
