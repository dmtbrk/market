package grpc

import (
	"context"
	"github.com/ortymid/market/jwt"
	"github.com/ortymid/market/market/user"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	Authorize(ctx context.Context, md metadata.MD) (*user.User, error)
	MetadataWithAuthorization(ctx context.Context, u *user.User) (metadata.MD, error)
}

type JWTAuthService struct {
	jwtService jwt.Service
}

func NewJWTAuthService(url string) *JWTAuthService {
	return &JWTAuthService{jwtService: jwt.Service{URL: url}}
}

func (s *JWTAuthService) Authorize(ctx context.Context, md metadata.MD) (*user.User, error) {
	values := md.Get("authorization")
	if values == nil || len(values) == 0 {
		// Anonymous call.
		return nil, nil
	}

	return s.jwtService.Authorize(ctx, values[0])
}

func (s *JWTAuthService) MetadataWithAuthorization(ctx context.Context, u *user.User) (metadata.MD, error) {
	panic("not implemented")
}

type UserIDAuthService struct{}

func (s *UserIDAuthService) Authorize(ctx context.Context, md metadata.MD) (*user.User, error) {
	values := md.Get("authorization")
	if values == nil || len(values) == 0 {
		// Anonymous call.
		return nil, nil
	}

	return &user.User{ID: values[0]}, nil
}

func (s *UserIDAuthService) MetadataWithAuthorization(ctx context.Context, u *user.User) (metadata.MD, error) {
	if u == nil {
		// Anonymous call.
		return metadata.New(map[string]string{}), nil
	}

	md := metadata.New(map[string]string{
		"authorization": u.ID,
	})
	return md, nil
}
