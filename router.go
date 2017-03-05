package ovhauth

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Secure installs the OVH authentication in a given gin-gonic router
// given a login url and a secret to encrypt user consumer key
func Secure(c *gin.Engine) *gin.RouterGroup {
	baseURL := os.Getenv("AUTH_LOGIN_URL")
	if len(baseURL) == 0 {
		logrus.Fatal("AUTH_LOGIN_URL is empty")
	}
	secret := os.Getenv("AUTH_SECRET")
	if len(secret) == 0 {
		logrus.Fatal("AUTH_SECRET is empty")
	}

	authModule := ovhAuthModule{
		baseURL: baseURL,
		secret:  secret,
	}

	c.GET("/auth/credential", authModule.GetCredential)
	c.GET("/auth/validate/:token", authModule.ValidateToken)

	authorized := c.Group("/")
	authorized.Use(jWTAuthMiddleware(secret))

	return authorized
}

type ovhAuthModule struct {
	baseURL string
	secret  string
}

// GetCredential calls the OVH API to get a validation URL and a consumer key.
// The consumer key is stored in memory and a token.
func (a *ovhAuthModule) GetCredential(c *gin.Context) {
	token := generateUUID()
	redirection := a.baseURL + "/?token=" + token

	ckValidationState, err := a.getConsumerKey(redirection)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	// Store consumerKey to retrieve it later
	keysMap.set(token, ckValidationState.ConsumerKey)

	c.JSON(200, gin.H{"url": ckValidationState.ValidationURL})
}

// ValidateToken retrieves a consumer key given a token
func (a *ovhAuthModule) ValidateToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		HTTPError(c, 400, errors.New("Token invalid"), nil)
		return
	}

	// Retrieve the consumer key given the token
	consumerKey := keysMap.get(token)
	if consumerKey == "" {
		HTTPError(c, 400, errors.New("Token invalid"), nil)
		return
	}

	// Get me
	me, err := a.GetMe(consumerKey)
	if err != nil {
		HTTPError(c, 401, errors.New("OVH authentication failed"), err)
		return
	}

	authUser := AuthUser{
		Me:          *me,
		ConsumerKey: consumerKey,
	}

	authUserString, err := json.Marshal(authUser)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	auth, err := SignAuth(string(authUserString), a.secret)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	c.JSON(200, gin.H{
		"auth":  auth,
		"name":  me.Name,
		"email": me.Email,
	})
}

///

// HTTPError logs an error and returns an HTTP error given a status cod, erre
func HTTPError(c *gin.Context, status int, userErr error, err error) {
	logrus.WithError(err).WithField("status", status)

	c.JSON(status, gin.H{"error": userErr.Error()})
}

///

type Map struct {
	lock sync.RWMutex
	Map  map[string]string
}

func (this *Map) get(key string) (v string) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.Map[key]; ok {
		return v
	}
	return ""
}

func (this *Map) set(key string, value string) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	this.Map[key] = value
}

// KeysMap keeps the user consumerKeys in memory
var keysMap = Map{
	Map: make(map[string]string),
}
