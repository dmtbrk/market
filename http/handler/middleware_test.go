package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTMiddleware(t *testing.T) {
	t.Run("Should populate context with user id", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OTAiLCJuYW1lIjoiSm9obiBEb2UiLCJpYXQiOjE1MTYyMzkwMjJ9.PKsVspL7mXlTfn55IQXLonedrX2XOaOrQ0gjOdDI0ek"
		alg := "HS256"
		secret := []byte("secret")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id, ok := r.Context().Value(KeyUserID).(string); !ok {
				t.Errorf("user id not in request context: got %q", id)
			}
			w.WriteHeader(http.StatusTeapot)
		})

		w := httptest.NewRecorder()
		m := JWTMiddleware(h, alg, secret)
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusTeapot {
			t.Errorf("unexpected status: got %q\nbody: %q", resp.Status, w.Body.Bytes())
		}
	})

	t.Run("Should do nothing on the absence of Authorization header", func(t *testing.T) {
		alg := "HS256"
		secret := []byte("secret")

		r := httptest.NewRequest(http.MethodGet, "/", nil)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(KeyUserID)
			if id != nil {
				t.Errorf("user id in request context: got %q", id)
			}
			w.WriteHeader(http.StatusTeapot)
		})

		w := httptest.NewRecorder()
		m := JWTMiddleware(h, alg, secret)
		m.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusTeapot {
			t.Errorf("unexpected status: got %q\nbody: %q", resp.Status, w.Body.Bytes())
		}
	})

	t.Run("Should responde with error on absense of user id", func(t *testing.T) {
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o"
		alg := "HS256"
		secret := []byte("secret")

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Errorf("request handling not stopped")
		})

		w := httptest.NewRecorder()
		m := JWTMiddleware(h, alg, secret)
		m.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("unexpected status: got %q\nbody: %q", resp.Status, w.Body.Bytes())
		}
	})
}
