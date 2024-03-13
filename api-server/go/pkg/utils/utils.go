package utils

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"github.com/caellach/shorturl/api-server/go/pkg/env"
)

func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func SplitAtUpperCase(s string) []string {
	var words []string
	var wordStart int
	for i, r := range s {
		if unicode.IsUpper(r) && i != 0 {
			words = append(words, s[wordStart:i])
			wordStart = i
		}
	}
	words = append(words, s[wordStart:])
	return words
}

func AddQueryParams(baseUrl string, params map[string]string) string {
	// check if the url already has query parameters
	baseUrlSplit := strings.Split(baseUrl, "?")
	if len(baseUrlSplit) > 1 {
		if len(params) > 0 {
			baseUrl += "&"
		}
	} else {
		baseUrl += "?"
	}

	// Add the query parameters to the url
	for key, value := range params {
		baseUrl += key + "=" + url.QueryEscape(value) + "&"
	}

	// Remove the trailing &
	return baseUrl[:len(baseUrl)-1]
}

const randomCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const randomCharsetLength = len(randomCharset)

func GenerateRandomString(length int) string {
	// Create a random string of the specified length
	b := make([]byte, length)
	for i := range b {
		b[i] = randomCharset[rand.Intn(randomCharsetLength)]
	}
	return string(b)
}

func GenerateErrorMessage(message string, err error) *fiber.Map {
	m := fiber.Map{
		"error": message,
	}
	if env.Config.Debug {
		if err != nil {
			m["error_message"] = err.Error()
		}
	}

	return &m
}

func LogError(logger fasthttp.Logger, response *fiber.Map) {
	m, err := json.Marshal(response)
	if err != nil {
		logger.Printf("error marshalling error message: %s", err)
		return
	}
	mStr := string(m)
	logger.Printf(mStr)
}

// LogServerError logs an error message and the error to the logger
// only for server errors where there is no request context (*fiber.Ctx)
func LogServerError(logger fasthttp.Logger, message string, err error) {
	response := GenerateErrorMessage(message, err)
	LogError(logger, response)
}

// GenerateJsonErrorMessage generates a json error message and logs it
// for server errors where there is a request context (*fiber.Ctx)
func GenerateJsonErrorMessage(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := GenerateErrorMessage(message, err)
	LogError(c.Context().Logger(), response)
	return c.Status(statusCode).JSON(response)
}

func getDefaultPort(protocol string) string {
	if protocol == "https" {
		return "443"
	}
	return "80"
}

func GetHost(c *fiber.Ctx) string {
	host := c.Hostname()
	port := c.Get("X-Forwarded-Port")
	if port == "" {
		rhost, rport, err := net.SplitHostPort(host)
		if err != nil {
			port = ""
		} else {
			host = rhost
			port = rport
		}
	}
	protocol := c.Protocol()

	if port != "" && port != getDefaultPort(protocol) {
		port = ":" + port
	} else {
		port = ""
	}

	return protocol + "://" + host + port
}

func GetRedirectUri(c *fiber.Ctx) string {
	return GetHost(c) + "/api/auth/callback"
}
