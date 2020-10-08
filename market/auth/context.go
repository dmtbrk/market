package auth

import (
	"context"
	"errors"
	"github.com/ortymid/market/market/user"
)

type contextKey struct{}

var (
	ContextKeyUser  contextKey
	ContextKeyToken contextKey
)

func NewContextWithUser(ctx context.Context, u *user.User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, u)
}

func UserFromContext(ctx context.Context) (*user.User, error) {
	u, ok := ctx.Value(ContextKeyUser).(*user.User)
	if !ok {
		return u, errors.New("user not provided")
	}
	return u, nil
}
