package ovhauth

import "github.com/ovh/go-ovh/ovh"

func (a *ovhAuthModule) getConsumerKey(redirection string) (*ovh.CkValidationState, error) {
	ovhClient, err := ovh.NewDefaultClient()
	if err != nil {
		return nil, err
	}

	ckRequest := ovhClient.NewCkRequest()
	ckRequest.AddRule("GET", "/me")

	ckValidationState, err := ckRequest.Do()
	if err != nil {
		return nil, err
	}

	return ckValidationState, nil
}

// -------------------------

// Me represents an OVH user
type Me struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// OvhGetMe calls the OVH API /me endpoint given an OVH config and a consumer key and
// unmarshals the JSON response in the value pointed to by result
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
