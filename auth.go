package ovhauth

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Register installs OVH authentication in a gin router
func Register(c *gin.Engine) error {
	baseURL := os.Getenv("AUTH_BASE_URL")
	if len(baseURL) == 0 {
		logrus.Fatal("AUTH_BASE_URL is empty")
	}

	authModule := ovhAuthModule{
		baseURL: baseURL,
	}

	authorized := c.Group("/api")
	authorized.Use(authRequired())

	c.GET("/auth/credential", authModule.GetCredential)
	c.GET("/auth/validate/:token", authModule.ValidateToken)

	return nil
}

type ovhAuthModule struct {
	baseURL string
}

// GetCredential calls the OVH API to get a validation URL and a consumer key
// The consumer key is stored in memory and a token
func (a *ovhAuthModule) GetCredential(c *gin.Context) {
	token := GenerateUUID()
	redirection := a.baseURL + "/?token=" + token

	ckValidationState, err := a.getConsumerKey(redirection)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	// Store consumerKey to retrieve it later
	KeysMap.Set(token, ckValidationState.ConsumerKey)

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
	consumerKey := KeysMap.Get(token)
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

	auth, err := Encrypt(CryptoKey, string(authUserString))
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	c.JSON(200, gin.H{
		"auth":  auth,
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

// get the map's value
func (this *Map) Get(key string) (v string) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.Map[key]; ok {
		return v
	}
	return ""
}

func (this *Map) Set(key string, value string) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	this.Map[key] = value
}

// KeysMap keeps the user consumerKeys in memory
var KeysMap = Map{
	Map: make(map[string]string),
}
