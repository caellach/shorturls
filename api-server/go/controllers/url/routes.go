package url

import (
	"errors"
	"fmt"
	"time"

	"github.com/caellach/shorturl/api-server/go/pkg/middleware"
	"github.com/caellach/shorturl/api-server/go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/caellach/shorturl/api-server/go/pkg/wordlist"
)

func deleteUrl(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AuthUser)

	id := c.Params("id")

	numberId, err := getNumberIdFromWordId(id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid id", err)
	}

	userObjectId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to convert user id to object id", err)
	}

	result, err := shorturlsCollection.UpdateOne(c.Context(), bson.M{"id": numberId, "userId": userObjectId, "deleted": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"deleted": primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))}})
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to delete url", err)
	}

	if result.ModifiedCount == 0 {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusNotFound, "url not found", fmt.Errorf("url not found: %s", id))
	}

	return c.SendStatus(fiber.StatusOK)
}

// Gets all urls for an authenticated user
func getUrls(c *fiber.Ctx) error {
	user := c.Locals("user").(middleware.AuthUser)

	userObjectId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to convert user id to object id", err)
	}

	cursor, err := shorturlsCollection.Find(c.Context(), bson.M{"userId": userObjectId, "deleted": nil})
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to find urls", err)
	}

	var urlsList []UrlDocument = make([]UrlDocument, 0)
	if err = cursor.All(c.Context(), &urlsList); err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusInternalServerError, "failed to retrieve urls", err)
	}

	// replace the ids in the list with the shorturl words
	// we don't want to expose the ids because then they can map the numbers to the words
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
func getUrlById(c *fiber.Ctx) error {
	id := c.Params("id")

	numberId, err := getNumberIdFromWordId(id)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "invalid id", err)
	}

	var url UrlDocument
	err = shorturlsCollection.FindOne(c.Context(), bson.M{"id": numberId, "deleted": bson.M{"$exists": false}}).Decode(&url)
	if err != nil {
		return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get url", err)
	}

	return c.Redirect(url.Url)
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
		Id:      urlId,
		UserId:  user.Id,
		Url:     putUrlRequest.Url,
		Created: primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	// the document should already exist because generateValidShortUrlId will only return a unique id if it can insert it into the collection
	// now we just need to update the document with the correct data
	_, err := shorturlsCollection.UpdateOne(c.Context(), bson.M{"id": urlId}, bson.M{"$set": url})
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

	return c.JSON(url)
}
