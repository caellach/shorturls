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
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	Email       string             `json:"email" bson:"email"`
	Username    string             `json:"username" bson:"username"`
	DisplayName string             `json:"displayName" bson:"displayName"`
	Locale      string             `json:"locale" bson:"locale"`
	MFAEnabled  bool               `json:"mfaEnabled" bson:"mfaEnabled"`
	Verified    bool               `json:"verified" bson:"verified"`
}

type UpdatedUserDocument struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	DisplayName string             `json:"displayName" bson:"displayName"`
	Avatar      string             `json:"avatar" bson:"avatar"`
	Email       string             `json:"email" bson:"email"`
	Locale      string             `json:"locale" bson:"locale"`
	MfaEnabled  bool               `json:"mfaEnabled" bson:"mfaEnabled"`
	Verified    bool               `json:"verified" bson:"verified"`
	Providers   []ProviderDocument `json:"providers" bson:"providers"`
}

type DiscordTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUserInfo struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"global_name"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Locale      string `json:"locale"`
	MfaEnabled  bool   `json:"mfa_enabled"`
	Verified    bool   `json:"verified"`
}
