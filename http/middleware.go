package http

import (
	"errors"
	"github.com/ortymid/market/market/auth"
	"net/http"
	"strings"
)

func AuthMiddleware(s AuthService, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtain token.
		//token, err := tokenFromHeader(r)
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//}

		// Create Authorize request.
		//ctx := auth.NewContextWithToken(context.Background(), token)

		// Make Authorize request.
		user, err := s.Authorize(r.Context(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		//
		r = r.WithContext(auth.NewContextWithUser(r.Context(), user))

		h.ServeHTTP(w, r)
	})
}

func tokenFromHeader(r *http.Request) (string, error) {
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
