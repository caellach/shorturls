package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/caellach/shorturl/api-server/go/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/caellach/shorturl/api-server/go/pkg/wordlist"
)

func deleteUrl(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AuthUser)

	id := c.Params("id")

	numberId, err := getNumberIdFromWordId(id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid id", err)
	}

	// mongo transaction
	session, err := mongoClient.StartSession()
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to start session", err)
	}
	defer session.EndSession(c.Context())

	// start the transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		// delete the url
		result, err := shorturlsCollection.UpdateOne(c.Context(), bson.M{"id": numberId, "user_id": user.Id, "deleted": bson.M{"$exists": false}},
			bson.M{"$set": bson.M{"deleted": primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))}})
		if err != nil {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to delete url", err)
		}

		if result.ModifiedCount == 0 {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusNotFound, "url not found", fmt.Errorf("url not found: %s", id))
		}

		// update the user metadata, decrement the active count, FindOneAndUpdate
		filter := bson.M{"user_id": user.Id}
		update := bson.M{"$inc": bson.M{"active_count": -1}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var userMetadata UserUrlMetadata
		err = metadataCollection.FindOneAndUpdate(c.Context(), filter, update, opts).Decode(&userMetadata)
		if err != nil {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to update user metadata", err)
		}
		return nil, nil
	}

	_, err = session.WithTransaction(c.Context(), callback)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to delete url", err)
	}

	go func() {
		// write to the websocket in the background
		// this is a fire and forget operation

		// generate the json message
		jsonMessage, err := json.Marshal(map[string]interface{}{"action": "deleted", "data": map[string]interface{}{"id": id}})
		if err != nil {
			log.Println("failed to marshal json message:", err)
			return
		}

		for _, conn := range websocketConnections[user.Id] {
			conn.WriteMessage(websocket.TextMessage, jsonMessage)
		}
	}()

	return c.SendStatus(fiber.StatusOK)
}

// Gets all urls for an authenticated user
func getUrls(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AuthUser)

	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}})
	cursor, err := shorturlsCollection.Find(c.Context(), bson.M{"user_id": user.Id, "deleted": nil}, opts)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to find urls", err)
	}

	var urlsList []UrlDocument = make([]UrlDocument, 0)
	if err = cursor.All(c.Context(), &urlsList); err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to retrieve urls", err)
	}

	// replace the ids in the list with the shorturl words
	// we don't want to expose the ids because users can map the numbers to the words
	badIds := make([]int, 0)
	for i, url := range urlsList {
		words := wordlist.GetWordsFromId(wordList, url.Id)
		if words == "" {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, fmt.Sprintf("invalid id: %s", url.Id), nil)
		}
		urlsList[i].Id = words
		urlsList[i].ShortUrl = generateFullShortUrl(c, words)
	}

	// remove the bad ids from the list
	// realistically this should never happen because the ids are generated from the wordlist
	// if someone changes the wordlist then the ids will change and some urls will be lost without conversion to the new wordlist
	for i := len(badIds) - 1; i >= 0; i-- {
		urlsList = append(urlsList[:badIds[i]], urlsList[badIds[i]+1:]...)
	}

	return c.JSON(urlsList)
}

// Redirects to the url with the given id
func redirectUrlById(c *fiber.Ctx) error {
	id := c.Params("id")

	numberId, err := getNumberIdFromWordId(id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid id", err)
	}

	var url UrlDocument
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	update := bson.M{"$set": bson.M{"last_used": primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))}, "$inc": bson.M{"use_count": 1}}
	err = shorturlsCollection.FindOneAndUpdate(c.Context(), bson.M{"id": numberId, "deleted": bson.M{"$exists": false}}, update, opts).Decode(&url)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get url", err)
	}

	go func() {
		url.Id = id // set the id to the word id

		jsonMessage, err := json.Marshal(map[string]interface{}{"action": "updated", "data": url})
		if err != nil {
			log.Println("failed to marshal json message:", err)
			return
		}

		for _, conn := range websocketConnections[url.UserId] {
			conn.WriteMessage(websocket.TextMessage, jsonMessage)
		}
	}()

	return c.Redirect(url.Url)
}

func getUserMetadata(c *fiber.Ctx) error {
	userMetadata, err := getMetadataForUser(c)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to get user metadata", err)
	}

	return c.JSON(userMetadata)
}

// Create a new url for the authenticated user
func putUrl(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AuthUser)

	var putUrlRequest PutUrlRequest = PutUrlRequest{}
	if err := c.BodyParser(&putUrlRequest); err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "content-type must be 'application/json'", err)
	}

	// validate the url
	if !isValidUrl(putUrlRequest.Url) {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, fmt.Sprintf("invalid url: %s", putUrlRequest.Url), fmt.Errorf("invalid url: %s", putUrlRequest.Url))
	}

	urlId := generateValidShortUrlId()
	if urlId == "" {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to generate url id", errors.New("failed to generate url id"))
	}

	url := UrlDocument{
		Id:       urlId,
		UserId:   user.Id,
		Url:      putUrlRequest.Url,
		Created:  primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond)),
		UseCount: 0,
		LastUsed: primitive.DateTime(0),
	}

	// mongo transaction
	session, err := mongoClient.StartSession()
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to start session", err)
	}
	defer session.EndSession(c.Context())

	// start the transaction
	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		// insert the url
		// the document should already exist because generateValidShortUrlId will only return a unique id if it can insert it into the collection
		// now we just need to update the document with the correct data
		_, err := shorturlsCollection.UpdateOne(c.Context(), bson.M{"id": urlId}, bson.M{"$set": url})
		if err != nil {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to insert url", err)
		}

		// update the user metadata, increment the active count, FindOneAndUpdate
		filter := bson.M{"user_id": user.Id}
		update := bson.M{"$inc": bson.M{"active_count": 1, "created_count": 1}, "$set": bson.M{"last_created": url.Created}}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var userMetadata UserUrlMetadata
		err = metadataCollection.FindOneAndUpdate(c.Context(), filter, update, opts).Decode(&userMetadata)
		if err != nil {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to update user metadata", err)
		}
		return nil, nil
	}
	_, err = session.WithTransaction(c.Context(), callback)

	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to insert url", err)
	}

	// convert the id to the shorturl
	words := wordlist.GetWordsFromId(wordList, urlId)
	if words == "" {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, fmt.Sprintf("invalid id: %s", urlId), nil)
	}
	url.Id = words
	url.ShortUrl = generateFullShortUrl(c, words)

	go func() {
		jsonMessage, err := json.Marshal(map[string]interface{}{"action": "created", "data": url})
		if err != nil {
			log.Println("failed to marshal json message:", err)
			return
		}

		for _, conn := range websocketConnections[url.UserId] {
			conn.WriteMessage(websocket.TextMessage, jsonMessage)
		}
	}()

	return c.JSON(url)
}

func getMetadataForUser(c *fiber.Ctx) (*UserUrlMetadata, error) {
	user := c.Locals("user").(middleware.AuthUser)

	filter := bson.M{"user_id": user.Id}

	var defaultUserUrlMetadata UserUrlMetadata = UserUrlMetadata{
		UserId:       user.Id,
		ActiveCount:  0,
		CreatedCount: 0,
		LastCreated:  primitive.DateTime(0),
	}

	// Define the update
	update := bson.M{
		"$setOnInsert": defaultUserUrlMetadata,
	}

	// Define the options
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	// Perform the FindOneAndUpdate operation
	var result *UserUrlMetadata
	err := metadataCollection.FindOneAndUpdate(c.Context(), filter, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func forceDisconnectWebsocket(wsContext *websocket.Conn, userId string) {
	log.Printf("WebSocket force disconnected: %s", wsContext.RemoteAddr().String())
	// kill the socket
	wsContext.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	if userId == "" {
		return
	}

	for i, conn := range websocketConnections[userId] {
		if conn == wsContext {
			websocketConnections[userId] = append(websocketConnections[userId][:i], websocketConnections[userId][i+1:]...)
		}
	}
}

var MAX_WS_MSG_SIZE = 1024

func urlWs(wsContext *websocket.Conn) {
	// WebSocket connected
	authorized := false
	userId := ""
	log.Printf("WebSocket connected: %s", wsContext.RemoteAddr().String())

	// Listen for messages infinitely

	for {
		msgType, msg, err := wsContext.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			forceDisconnectWebsocket(wsContext, userId)
			break
		}

		if !authorized {
			if msgType == websocket.TextMessage {
				if len(msg) == 0 || len(msg) > MAX_WS_MSG_SIZE {
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// decode the message
				var message map[string]interface{}
				err := json.Unmarshal(msg, &message)
				if err != nil {
					log.Println("failed to unmarshal json message:", err)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// check the action
				action, ok := message["action"].(string)
				if !ok {
					log.Println("invalid message: action not found")
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				if action == "auth" {
					// check the token
					token, ok := message["token"].(string)
					if !ok {
						log.Println("invalid message: token not found")
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// check the token, it should be signed by the server
					decodedToken, err := middleware.ValidateToken(token)
					if err != nil {
						log.Println("failed to verify token:", err)
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// check the user id
					userId, ok = decodedToken.Claims.(jwt.MapClaims)["sub"].(string)
					if !ok {
						log.Println("invalid token: user id not found")
						forceDisconnectWebsocket(wsContext, userId)
						break
					}

					// register the websocket
					websocketConnections[userId] = append(websocketConnections[userId], wsContext)
					authorized = true
					// send the auth response
					wsContext.WriteMessage(websocket.TextMessage, []byte(`{"action":"auth"}`))
					continue
				}
			}
			forceDisconnectWebsocket(wsContext, userId)
			break
		} else {
			if msgType == websocket.PingMessage {
				if string(msg) == "ping" {
					wsContext.WriteMessage(websocket.PongMessage, []byte("pong"))
				}
			} else if msgType == websocket.TextMessage {
				if len(msg) == 0 || len(msg) > MAX_WS_MSG_SIZE {
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				var message map[string]interface{}
				err := json.Unmarshal(msg, &message)
				if err != nil {
					log.Println("failed to unmarshal json message:", err)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				// check the action
				action, ok := message["action"].(string)
				if !ok {
					log.Println("invalid message: action not found")
					forceDisconnectWebsocket(wsContext, userId)
					break
				}

				if action == "ping" {
					wsContext.WriteMessage(websocket.TextMessage, []byte(`{"action":"pong"}`))
					continue
				} else {
					log.Println("invalid action:", action)
					forceDisconnectWebsocket(wsContext, userId)
					break
				}
			}
		}
	}
}
