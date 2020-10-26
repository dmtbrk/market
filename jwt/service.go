package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/ortymid/market/market/user"
	"io/ioutil"
	"net/http"
)

type Service struct {
	URL string

	cacheSecret interface{}
}

func (s *Service) Authorize(ctx context.Context, token interface{}) (*user.User, error) {
	tokenString, ok := token.(string)
	if !ok {
		return nil, nil
	}

	claims, err := Parse(tokenString, s.SecretContext(ctx))
	if err != nil {
		return nil, err
	}

	return &user.User{ID: claims.UserID}, nil
}

func (s *Service) SecretContext(ctx context.Context) func() (interface{}, error) {
	return func() (interface{}, error) {
		if s.cacheSecret == nil {
			secret, err := s.fetchSecret(ctx)
			if err != nil {
				return nil, err
			}
			s.cacheSecret = secret
		}

		return s.cacheSecret, nil
	}
}

func (s *Service) fetchSecret(ctx context.Context) (interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("reading secret response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching secret: %s", body)
	}

	key := &rsa.PublicKey{}
	err = json.Unmarshal(body, key)

	return key, err
}
