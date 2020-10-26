package auth

import (
	"context"
	"github.com/ortymid/market/market/user"
)

type Service interface {
	Authorize(ctx context.Context, token interface{}) (*user.User, error)
}
