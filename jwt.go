package ovhauth

import "github.com/dgrijalva/jwt-go"

// SignAuth signs auth data using a secret
func SignAuth(auth string, secret string) (string, error) {
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth": auth,
	})

	jwtokenString, err := jwtoken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return jwtokenString, nil
}
