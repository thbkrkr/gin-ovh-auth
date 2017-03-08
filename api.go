package ovhauth

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/ovh/go-ovh/ovh"
)

type ovhAuthModule struct {
	secret string
}

// cKeyCache keeps the user consumerKeys in memory
var cKeyCache = CKeyCache{
	Map: make(map[string]string),
}

// GetCredential calls the OVH API to get a validation URL and a consumer key.
// The consumer key is stored in memory and a token.
func (a *ovhAuthModule) GetCredential(c *gin.Context) {
	token := generateUUID()

	redirect, ok := c.GetQuery("redirect")
	if !ok {
		err := errors.New("Redirect URL not found")
		HTTPError(c, 400, err, err)
		return
	}
	redirection, err := url.QueryUnescape(redirect)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	rules := map[string]string{
		"GET":    "/*",
		"POST":   "/*",
		"PUT":    "/*",
		"DELETE": "/*",
	}

	ckValidationState, err := a.getConsumerKey(rules, redirection)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	// Store consumerKey to retrieve it later
	cKeyCache.set(token, ckValidationState.ConsumerKey)

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
	consumerKey, ko := cKeyCache.get(token)
	if ko {
		HTTPError(c, 400, errors.New("Token invalid"), nil)
		return
	}

	cKeyCache.delete(token)

	// Get me
	me, err := a.GetMe(consumerKey)
	if err != nil {
		HTTPError(c, 401, errors.New("Authentication failed"), err)
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

// HTTPError logs an error and returns an HTTP error given a status cod, erre
func HTTPError(c *gin.Context, status int, userErr error, err error) {
	logrus.WithError(err).WithField("status", status)
	c.JSON(status, gin.H{"error": userErr.Error()})
}

func (a *ovhAuthModule) getConsumerKey(rules map[string]string, redirection string) (*ovh.CkValidationState, error) {
	ovhClient, err := ovh.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	// TODO be able to parametrize rules
	ckRequest := ovhClient.NewCkRequestWithRedirection(redirection)
	ckRequest.AddRule("GET", "/*")
	ckRequest.AddRule("POST", "/*")
	ckRequest.AddRule("PUT", "/*")
	ckRequest.AddRule("DELETE", "/*")

	ckValidationState, err := ckRequest.Do()
	if err != nil {
		return nil, err
	}

	return ckValidationState, nil
}

// Me represents an OVH user
type Me struct {
	ID    string `json:"nichandle"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetMe calls the OVH API /me endpoint given a consumer key and
// return the corresponding Me
func (a *ovhAuthModule) GetMe(consumerKey string) (*Me, error) {
	ovhClient, err := ovh.NewDefaultClient()
	ovhClient.ConsumerKey = consumerKey
	if err != nil {
		return nil, err
	}

	var me Me
	err = ovhClient.Get("/me", &me)
	if err != nil {
		return nil, err
	}

	return &me, nil
}
