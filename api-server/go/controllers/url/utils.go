package url

import (
	"context"
	"fmt"
	"math/rand"
	netUrl "net/url"
	"strings"

	"github.com/caellach/shorturl/api-server/go/pkg/utils"
	"github.com/caellach/shorturl/api-server/go/pkg/wordlist"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidUrl(url string) bool {
	// must start with http:// or https://
	if url[:7] != "http://" && url[:8] != "https://" {
		return false
	}

	_, err := netUrl.ParseRequestURI(url)
	//nolint
	return err == nil
}

// returns a number between 0 and 8191 inclusive as a formatted string, e.g. 0000
func generateRandomUrlIdPart() string {
	// 0-8191 inclusive
	num := rand.Intn(wordList.Count)
	return fmt.Sprintf("%04d", num)
}

func generateValidShortUrlId() string {
	// will generate a string of numbers in the format 000-000-000-000 and check if it exists in the database
	// if it does, it will generate a new one and check again
	// this will continue until a unique id is found

	// minimal UrlDocument struct to hold the id
	urlDocument := UrlDocument{
		Id: "",
	}

	isUnique := false
	for !isUnique {
		id := []string{generateRandomUrlIdPart(), generateRandomUrlIdPart(), generateRandomUrlIdPart(), generateRandomUrlIdPart()}
		// check if the id exists in the database by trying to insert a document only containing the id
		// if it fails with a duplicate key error, continue the loop
		// if it succeeds, set isUnique to true
		urlDocument.Id = strings.Join(id, "-")
		_, err := shorturlsCollection.InsertOne(context.Background(), urlDocument)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				// log
				fmt.Println("duplicate key error (" + urlDocument.Id + "), generating new id")
				//continue
				return ""
			}
			fmt.Println("error inserting url:", err)
			return ""
		} else {
			isUnique = true
		}
	}
	return urlDocument.Id
}

func getNumberIdFromWordId(wordId string) (string, error) {
	// split the words by uppercase letters
	words := utils.SplitAtUpperCase(wordId)
	if len(words) != 4 {
		return "", fmt.Errorf("invalid id, bad length: %s", wordId)
	}

	// convert the words to ids
	numberId := wordlist.GetIdFromWords(wordList, words)
	return numberId, nil
}

func generateFullShortUrl(c *fiber.Ctx, wordId string) string {
	// get the short url from the number id
	return fmt.Sprintf("%s/u/%s", c.BaseURL(), wordId)
}
