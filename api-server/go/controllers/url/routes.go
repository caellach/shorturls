package url

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/caellach/shorturl/api-server/go/pkg/utils"
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
		result, err := shorturlsCollection.UpdateOne(c.Context(), bson.M{"id": numberId, "userId": user.Id, "deleted": bson.M{"$exists": false}},
			bson.M{"$set": bson.M{"deleted": primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))}})
		if err != nil {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to delete url", err)
		}

		if result.ModifiedCount == 0 {
			return nil, utils.GenerateJsonErrorMessage(c, fiber.StatusNotFound, "url not found", fmt.Errorf("url not found: %s", id))
		}

		// update the user metadata, decrement the active count, FindOneAndUpdate
		filter := bson.M{"userId": user.Id}
		update := bson.M{"$inc": bson.M{"activeCount": -1}}
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
	cursor, err := shorturlsCollection.Find(c.Context(), bson.M{"userId": user.Id, "deleted": nil}, opts)
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

	isEmbed := false
	userAgent := strings.ToLower(c.Get("User-Agent"))
	for _, embedUserAgent := range *embedUserAgents {
		if strings.Contains(userAgent, embedUserAgent) {
			isEmbed = true
			break
		}
	}

	var url UrlDocument
	if isEmbed {
		// get the url data without updating
		err = shorturlsCollection.FindOne(c.Context(), bson.M{"id": numberId, "deleted": bson.M{"$exists": false}}).Decode(&url)
		if err != nil {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get url", err)
		}

		if !url.OgpDataId.IsZero() {
			return c.Redirect(fmt.Sprintf("/u/f/%s", url.OgpDataId.Hex()))
		} else {
			return c.Redirect(url.Url)
		}
	} else {
		// log user agent
		log.Println("user agent:", userAgent)

		// get the url data and update the lastUsed and useCount
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		update := bson.M{"$set": bson.M{"lastUsed": primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))}, "$inc": bson.M{"useCount": 1}}
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
		filter := bson.M{"userId": user.Id}
		update := bson.M{"$inc": bson.M{"activeCount": 1, "createdCount": 1}, "$set": bson.M{"lastCreated": url.Created}}
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

	filter := bson.M{"userId": user.Id}

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

func fakeOGPResult(c *fiber.Ctx) error {
	// return an html page with the ogp data
	// this is a fake response to test the ogp data

	// get id
	id := c.Params("id")
	log.Println("id:", id)

	// convert the id to the object id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid id", err)
	}

	/*ogpData := OgpData{
		Id:          primitive.NewObjectID(),
		SiteName:    "IMDb",
		Title:       "The Rock",
		Type:        "video.movie",
		Url:         "https://o7.rip/u/DespiseDevelopmentLetterBang",
		Image:       "https://static.miraheze.org/greatcharacterswiki/thumb/8/86/D75zqo-a1d18cbe-acbc-4f91-9009-25b63e297eee.jpg/330px-D75zqo-a1d18cbe-acbc-4f91-9009-25b63e297eee.jpg",
		Description: "The Rock is a 1996 American action thriller film directed by Michael Bay, produced by Don Simpson and Jerry Bruckheimer, and written by David Weisberg and Douglas S. Cook.",
	}*/
	// get stored ogp data from mongo
	var ogpData OgpData
	err = ogpDataCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&ogpData)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get ogp data", err)
	}

	c.Type("html")
	return c.SendString(
		`<html prefix="og: http://ogp.me/ns#">
	<head>
		<meta name="twitter:card" content="summary_large_image">
		<meta property="og:site_name" content="` + ogpData.SiteName + `" />
		<meta property="og:title" content="` + ogpData.Title + `" />
		<meta property="og:type" content="` + ogpData.Type + `" />
		<meta property="og:url" content="` + ogpData.Url + `" />
		<meta property="og:image" content="` + ogpData.Image + `" />
		<meta property="og:description" content="` + ogpData.Description + `" />
	</head><body>` + utils.GenerateRandomString(32) + `</body>
</html>`)
}
