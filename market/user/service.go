package user

import "context"

type Service interface {
	Get(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, r UpdateRequest) (User, error)
}
