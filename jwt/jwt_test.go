package jwt

import (
	"reflect"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestParse(t *testing.T) {
	// key, _ := rsa.GenerateKey(rand.Reader, 2048)
	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, wantClaims)
	// tokenString, err := token.SignedString(key)
	// if err != nil {
	// 	t.Errorf("unexpected error: %v", err)
	// 	return
	// }
	t.Run("Should parse token with string user id", func(t *testing.T) {
		alg := "HS256"
		secret := []byte("secret")
		// { id: "123456789" }
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OSJ9.XMgTcgoDgqbxLvDQG6h6EjthqMhb4KBFY-nnYzAq-og"
		wantClaims := &Claims{
			UserID: "123456789",
		}

		got, err := Parse(token, alg, secret)
		if err != nil {
			t.Errorf("Parse() unexpected error: %v", err)
			return
		}
		if !reflect.DeepEqual(got, wantClaims) {
			t.Errorf("Parse() = %#v, want %#v", got, wantClaims)
		}
		t.Logf("got: %#v", got)
	})

	t.Run("Should parse token with int user id", func(t *testing.T) {
		alg := "HS256"
		secret := []byte("secret")
		// { id: 123456789 }
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2Nzg5fQ.xNQJDKQz08Zm1-MTTQNLj34CFu51WAm_WpRwnF9eAFU"
		wantClaims := &Claims{
			UserID: "123456789",
		}

		got, err := Parse(token, alg, secret)
		if err != nil {
			t.Errorf("Parse() unexpected error: %v", err)
			return
		}
		if !reflect.DeepEqual(got, wantClaims) {
			t.Errorf("Parse() = %#v, want %#v", got, wantClaims)
		}
		t.Logf("got: %#v", got)
	})

	t.Run("Should parse token with standart claims", func(t *testing.T) {
		alg := "HS256"
		secret := []byte("secret")
		// { id: "123456789", iat: 1516239022 }
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OSIsImlhdCI6MTUxNjIzOTAyMn0.SQiePc7E2ZJp8PjEU2AjcJ7JsD6_sT_gEUeg3XWHQt0"
		wantClaims := &Claims{
			UserID: "123456789",
			StandardClaims: jwt.StandardClaims{
				IssuedAt: 1516239022,
			},
		}

		got, err := Parse(token, alg, secret)
		if err != nil {
			t.Errorf("Parse() unexpected error: %v", err)
			return
		}
		if !reflect.DeepEqual(got, wantClaims) {
			t.Errorf("Parse() = %#v, want %#v", got, wantClaims)
		}
		t.Logf("got: %#v", got)
	})
}
