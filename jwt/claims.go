package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}
