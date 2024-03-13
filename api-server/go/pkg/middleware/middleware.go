package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"github.com/caellach/shorturl/api-server/go/pkg/config"
	"github.com/caellach/shorturl/api-server/go/pkg/utils"
)

// AuthRequired is a middleware that checks if the user is authenticated
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for a session or a JWT or some other form of authentication
		// For example, let's say we're checking for a JWT in the Authorization header
		authorization := c.Get("Authorization")

		token := strings.Split(authorization, " ")
		if strings.ToLower(token[0]) != "bearer" || len(token) != 2 {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to get token from request", errors.New("token is empty"))
		}

		parsedToken, err := ValidateToken(token[1])
		if err != nil {
			return utils.GenerateJsonErrorMessage(c, fiber.StatusBadRequest, "failed to validate token", err)
		}

		/*"sub":          updatedDocument.Id.Hex(),
		"username":     updatedDocument.Username,
		"avatar":       updatedDocument.Avatar,
		"provider":     authState.Provider,
		"provider_sub": providerId,
		"exp":          expiresAt,*/
		// If the token is valid, store the user information in Locals
		user := AuthUser{
			Id:          parsedToken.Claims.(jwt.MapClaims)["sub"].(string),
			Username:    parsedToken.Claims.(jwt.MapClaims)["username"].(string),
			Avatar:      parsedToken.Claims.(jwt.MapClaims)["avatar"].(string),
			Provider:    parsedToken.Claims.(jwt.MapClaims)["provider"].(string),
			ProviderSub: parsedToken.Claims.(jwt.MapClaims)["provider_sub"].(string),
			ExpiresAt:   int64(parsedToken.Claims.(jwt.MapClaims)["exp"].(float64)),
		}
		c.Locals("user", user)

		// Call the next handler
		return c.Next()
	}
}

func ValidateToken(token string) (*jwt.Token, error) {
	mySigningKeyFunc := func(token *jwt.Token) (interface{}, error) {
		// Check if the token is valid
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Need more checks

		// Return the key used to sign the token
		return []byte(config.ServerConfig.Token.Secret), nil
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, mySigningKeyFunc)
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return parsedToken, nil
}
