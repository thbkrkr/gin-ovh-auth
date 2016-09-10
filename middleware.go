package ovhauth

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	Me          Me     `json:"me"`
	ConsumerKey string `json:"consumerKey"`
}

var (
	authUserHeaderName = "X-Ovh-Auth"
)

// Middleware to set extract auth in the HTTP gin context
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := extractAuth(c)
		if err != nil {
			c.AbortWithStatus(401)
		} else {
			c.Set("auth", auth)
			c.Next()
		}
	}
}

// extractAuth retrieves an authentication user from the HTTP header X-Ovh-Auth
func extractAuth(c *gin.Context) (*AuthUser, error) {
	encryptedMe := c.Request.Header.Get(authUserHeaderName)

	jsonMe, err := Decrypt(CryptoKey, encryptedMe)
	if err != nil {
		return nil, err
	}

	authUser := AuthUser{}
	json.Unmarshal([]byte(jsonMe), &authUser)

	if authUser.Me.Email == "" || authUser.ConsumerKey == "" {
		return nil, errors.New("Empty auth user")
	}

	return &authUser, nil
}

func GetAuthUser(c *gin.Context) *AuthUser {
	return c.MustGet("auth").(*AuthUser)
}
