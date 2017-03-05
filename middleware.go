package ovhauth

import (
	"encoding/json"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	Me          Me     `json:"me"`
	ConsumerKey string `json:"consumerKey"`
}

func GetAuthUser(c *gin.Context) *AuthUser {
	return c.MustGet("auth").(*AuthUser)
}

var (
	authUserHeaderName = "X-Auth"
)

func jWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := c.Request.Header.Get("X-Auth")

		if rawToken == "" {
			c.AbortWithError(401, errors.New("Authorization failed"))
			return
		}

		token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		} else if token == nil {
			c.AbortWithError(401, errors.New("Authorization failed"))
		} else {

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				c.AbortWithError(401, errors.New("Authorization failed"))
			}

			sauthUser := claims["auth"].(string)
			var authUser AuthUser
			err := json.Unmarshal([]byte(sauthUser), &authUser)
			if err != nil {
				logrus.WithError(err).Error("Fail to unmarshal auth user")
				c.AbortWithError(401, errors.New("Authorization failed"))
				return
			}

			c.Set("auth", &authUser)
		}
	}
}
