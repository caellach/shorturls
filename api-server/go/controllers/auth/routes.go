package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caellach/shorturl/api-server/go/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func authProviderCallback(c *fiber.Ctx) error {
	// get the state from the request
	state := c.Query("state")
	if state == "" {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get state from request", errors.New("state is empty"))
	}

	// check if the state exists in the database
	var authState AuthStateDocument
	err := authStatesCollection.FindOne(c.Context(), map[string]string{
		"state": state,
	}).Decode(&authState)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to find state", err)
	}

	// get the code from the request
	code := c.Query("code")
	if code == "" {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get code from request", errors.New("code is empty"))
	}

	// get the token from the auth provider from the code

	username := ""
	avatar := ""
	providerId := ""
	email := ""
	if authState.Provider == "discord" {
		// pull username from the token
		tokenResponse, err := getTokenFromCodeDiscord(code)
		if err != nil {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to get token from code", err)
		}

		discordUserInfo, err := getDiscordMeUserInfo(tokenResponse.AccessToken)
		if err != nil {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to get user info from token", err)
		}

		providerId = discordUserInfo.Id
		username = discordUserInfo.Username
		avatar = discordUserInfo.Avatar
		email = discordUserInfo.Email
	} else {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid provider", errors.New("invalid provider"))
	}

	// create a new user in mongo if one doesn't exist; else, update the user
	// FindOneAndUpdate
	newDocument := UserDocument{ // to avoid an upsert conflict we create the Document without the providers
		Avatar:   avatar,
		Email:    email,
		Username: username,
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.M{
		"providers": bson.M{
			"$elemMatch": bson.M{
				"name": authState.Provider,
				"id":   providerId,
			},
		},
	}
	update := bson.M{"$set": newDocument,
		"$addToSet": bson.M{
			"providers": bson.M{
				"$each": []bson.D{ // preserve order to prevent duplicates
					{
						{Key: "name", Value: authState.Provider},
						{Key: "id", Value: providerId},
					},
				},
			},
		},
	}

	var updatedDocument UpdatedUserDocument
	err = usersCollection.FindOneAndUpdate(c.Context(), filter, update, opts).Decode(&updatedDocument)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to create or update user", err)
	}

	// create a new JWT token
	expiresAt := time.Now().Add(time.Hour * 6).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":          updatedDocument.Id.Hex(),
		"username":     updatedDocument.Username,
		"avatar":       updatedDocument.Avatar,
		"provider":     authState.Provider,
		"provider_sub": providerId,
		"exp":          expiresAt,
	})

	// sign the token
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to sign token", err)
	}

	// return the token
	/*return c.JSON(fiber.Map{
		"access_token":  tokenString,
		"expires_at":    expiresAt,
		"token_type":    "Bearer",
		"refresh_token": nil,
	})*/
	return c.Redirect(authState.Referer + "?a=" + tokenString)
}

func getTokenFromCodeDiscord(code string) (DiscordTokenResponse, error) {
	// get the token from the auth provider from the code
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:3000/api/auth/callback")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", discordConfig.ApiBaseUrl+"/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return DiscordTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discordConfig.ClientID, discordConfig.ClientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return DiscordTokenResponse{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DiscordTokenResponse{}, err
	}

	if strings.Contains(string(body), "error") {
		return DiscordTokenResponse{}, errors.New("failed to get token from code")
	}

	var discordTokenResponse DiscordTokenResponse
	err = json.Unmarshal(body, &discordTokenResponse)
	if err != nil {
		return DiscordTokenResponse{}, err
	}

	return discordTokenResponse, nil
}

func getDiscordMeUserInfo(token string) (DiscordUserInfo, error) {
	// get the user info from the token
	req, err := http.NewRequest("GET", discordConfig.ApiBaseUrl+"/users/@me", nil)
	if err != nil {
		return DiscordUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return DiscordUserInfo{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DiscordUserInfo{}, err
	}

	var discordUserInfo DiscordUserInfo
	err = json.Unmarshal(body, &discordUserInfo)
	if err != nil {
		return DiscordUserInfo{}, err
	}

	return discordUserInfo, nil
}

func getAuthProviderOAuthURL(c *fiber.Ctx) error {
	authProvider := c.Params("authProvider")

	// check if the auth provider is valid
	authProviderBaseUrl, exists := validAuthProviders[authProvider]
	if !exists {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid auth provider", errors.New("invalid auth provider"))
	}

	referer := c.Get("Referer")
	host := c.Get("Host")

	state := utils.GenerateRandomString(16)
	authProviderBaseUrl = utils.AddQueryParams(authProviderBaseUrl, map[string]string{
		"state":        state,
		"redirect_uri": "http://" + host + "/api/auth/callback",
	})

	// get request ip
	ip := c.IP()
	// get current date
	dateTime := time.Now()
	// store the state in Mongo
	authStateDocument := AuthStateDocument{
		State:      state,
		Ip:         ip,
		Provider:   authProvider,
		InsertedAt: primitive.NewDateTimeFromTime(dateTime),
		Referer:    referer,
	}
	_, err := authStatesCollection.InsertOne(c.Context(), authStateDocument)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to store state", err)
	}

	// return json with the auth provider url
	return c.JSON(fiber.Map{
		"url":      authProviderBaseUrl,
		"provider": authProvider,
	})
}
