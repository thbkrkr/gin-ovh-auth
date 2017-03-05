package ovhauth

import "github.com/ovh/go-ovh/ovh"

func (a *ovhAuthModule) getConsumerKey(redirection string) (*ovh.CkValidationState, error) {
	ovhClient, err := ovh.NewDefaultClient()
	if err != nil {
		return nil, err
	}

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
