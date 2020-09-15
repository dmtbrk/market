package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTMiddleware(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
	alg := "HS256"
	secret := "secret"

	t.Run("Should populate context with user id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Authorization", "Bearer "+token)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id, ok := r.Context().Value(KeyUserID).(string); !ok {
				t.Errorf("user id is not in request context: got %q", id)
			}
		})

		w := httptest.NewRecorder()
		m := JWTMiddleware(h, alg, secret)
		m.ServeHTTP(w, r)
	})

	t.Run("Should populate context with user id", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add("Authorization", "Bearer "+token)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id, ok := r.Context().Value(KeyUserID).(string); !ok {
				t.Errorf("user id is not in request context: got %q", id)
			}
		})

		w := httptest.NewRecorder()
		m := JWTMiddleware(h, alg, secret)
		m.ServeHTTP(w, r)
	})
}	
