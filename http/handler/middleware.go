package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ortymid/market/jwt"
)

// JWTMiddleware attaches a user ID obtained from the JWT to the request context.
// In case of invalid token the Forbidden response is written.
func JWTMiddleware(h http.Handler, alg string, secret interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenString(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if len(tokenString) == 0 { // no token is ok
			h.ServeHTTP(w, r)
			return
		}

		claims, err := jwt.Parse(tokenString, alg, secret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if len(claims.UserID) == 0 {
			http.Error(w, "user id is not in JWT claims", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyUserID, claims.UserID)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getTokenString looks for the JWT in the Authorization header.
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
		return "", errors.New("Authorization type is not Bearer")
	}

	token := authFields[1]
	return token, nil
}

func PaginationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		offset, err := strconv.Atoi(params.Get("offset"))
		if err != nil {
			http.Error(w, "valid offset query parameter required", http.StatusBadRequest)
			return
		}
		limit, err := strconv.Atoi(params.Get("limit"))
		if err != nil {
			http.Error(w, "valid limit query parameter required", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyPage, Page{
			Offset: offset,
			Limit:  limit,
		})
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type Page struct {
	Limit  int
	Offset int
}
