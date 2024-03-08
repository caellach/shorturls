package auth

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthStateDocument struct {
	State      string             `json:"state" bson:"state"`
	Ip         string             `json:"ip" bson:"ip"`
	InsertedAt primitive.DateTime `json:"insertedAt" bson:"insertedAt"`
	Provider   string             `json:"provider" bson:"provider"`
	Referer    string             `json:"referer" bson:"referer"`
}

type ProviderDocument struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type UserDocument struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Email    string             `json:"email" bson:"email"`
	Username string             `json:"username" bson:"username"`
}

type UpdatedUserDocument struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Email     string             `json:"email" bson:"email"`
	Username  string             `json:"username" bson:"username"`
	Providers []ProviderDocument `json:"providers" bson:"providers"`
}

type DiscordTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUserInfo struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}
