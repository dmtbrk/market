package product

import "context"

type Interface interface {
	List(ctx context.Context, r ListRequest) ([]*Product, error)
	Get(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, r CreateRequest) (*Product, error)
	Update(ctx context.Context, r UpdateRequest) (*Product, error)
	Delete(ctx context.Context, id string) (*Product, error)
}
