package ovhauth

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Secure installs the OVH authentication in a given gin-gonic router
// given a redirect url and a secret to encrypt user consumer key
func Secure(c *gin.Engine) *gin.RouterGroup {
	secret := os.Getenv("AUTH_SECRET")
	if len(secret) == 0 {
		logrus.Fatal("AUTH_SECRET is empty")
	}

	authModule := ovhAuthModule{
		secret: secret,
	}

	c.GET("/auth/credential", authModule.GetCredential)
	c.GET("/auth/validate/:token", authModule.ValidateToken)

	authorized := c.Group("/")
	authorized.Use(jWTAuthMiddleware(secret))

	return authorized
}
