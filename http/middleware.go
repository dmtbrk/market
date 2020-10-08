package http

import (
	"github.com/ortymid/market/market/auth"
	"net/http"
)

func AuthMiddleware(s AuthService, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.Authorize(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(auth.NewContextWithUser(r.Context(), user))

		h.ServeHTTP(w, r)
	})
}
