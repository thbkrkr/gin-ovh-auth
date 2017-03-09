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
	rules  map[string]string
}

// cache keeps the user consumerKeys in memory
var cache = cKeyCache{
	Map: map[string]string{},
}

// GetConsumerKey calls the OVH API to get a consumer key and
// an URL to redirect user in order to log in.
// The consumer key is stored in memory and a token is set is as
// query parameter in the redirect URL.
func (a *ovhAuthModule) GetConsumerKey(c *gin.Context) {
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

	if a.rules == nil {
		a.rules = map[string]string{}
	}
	if a.rules["GET"] == "" {
		a.rules["GET"] = "/*"
	}
	if a.rules["POST"] == "" {
		a.rules["POST"] = "/*"
	}
	if a.rules["PUT"] == "" {
		a.rules["PUT"] = "/*"
	}
	if a.rules["DELETE"] == "" {
		a.rules["DELETE"] = "/*"
	}

	ckValidationState, err := a.getConsumerKey(a.rules, redirection)
	if err != nil {
		HTTPError(c, 400, err, err)
		return
	}

	// Store consumerKey to retrieve it later
	cache.set(token, ckValidationState.ConsumerKey)

	c.JSON(200, gin.H{
		"url":   ckValidationState.ValidationURL,
		"token": token,
	})
}

// ValidateToken retrieves a consumer key given a token
func (a *ovhAuthModule) ValidateToken(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		HTTPError(c, 400, errors.New("Token invalid"), nil)
		return
	}

	// Retrieve the consumer key given the token
	consumerKey, ko := cache.get(token)
	if ko {
		HTTPError(c, 400, errors.New("Token invalid"), nil)
		return
	}

	cache.delete(token)

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

func (a *ovhAuthModule) getConsumerKey(rules map[string]string, redirection string) (*ovh.CkValidationState, error) {
	ovhClient, err := ovh.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	ckRequest := ovhClient.NewCkRequestWithRedirection(redirection)
	ckRequest.AddRule("GET", rules["GET"])
	ckRequest.AddRule("POST", rules["POST"])
	ckRequest.AddRule("PUT", rules["PUT"])
	ckRequest.AddRule("DELETE", rules["DELETE"])

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

// HTTPError logs an error and returns an HTTP error given a status code
func HTTPError(c *gin.Context, status int, retErr error, err error) {
	if err == nil {
		err = retErr
	}
	logrus.WithError(err).WithField("status", status)
	c.JSON(status, gin.H{"error": retErr.Error()})
}
