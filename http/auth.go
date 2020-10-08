package http

import (
	"errors"
	"github.com/ortymid/market/jwt"
	"github.com/ortymid/market/market/user"
	"net/http"
	"strings"
)

type AuthService interface {
	Authorize(r *http.Request) (*user.User, error)
}

type JWTAuthService struct {
	jwtService jwt.Service
}

func NewJWTAuthService(url string) *JWTAuthService {
	return &JWTAuthService{jwtService: jwt.Service{URL: url}}
}

func (s *JWTAuthService) Authorize(r *http.Request) (*user.User, error) {
	tokenString, err := getTokenString(r)
	if err != nil {
		return nil, err
	}

	if len(tokenString) == 0 {
		// Anonymous request.
		return nil, nil
	}

	return s.jwtService.Authorize(r.Context(), tokenString)
}

// getTokenString looks for JWT in the Authorization header.
// An empty token with nil error indicates the absence of Authorization header.
func getTokenString(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", nil // no Authorization header
	}

	authFields := strings.Fields(auth)
	if len(authFields) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	typ := authFields[0]
	if !strings.EqualFold(typ, "Bearer") {
		return "", errors.New("authorization type is not Bearer")
	}

	token := authFields[1]
	return token, nil
}
