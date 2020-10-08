package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type SecretFunc func() (interface{}, error)

func Parse(token string, secret SecretFunc) (Claims, error) {
	var claims Claims

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secret()
	})
	if err != nil {
		return claims, fmt.Errorf("parsing JWT: %w", err)
	}

	return claims, nil
}
