package product

import "context"

//go:generate mockgen -destination=../../mock/product_storage.go -package mock -mock_names=Storage=ProductStorage . Storage

type Storage interface {
	Finder
	Creater
	Updater
	Deleter
}

type Finder interface {
	Find(ctx context.Context, r FindRequest) ([]*Product, error)
	FindOne(ctx context.Context, id string) (*Product, error)
}

type Creater interface {
	Create(ctx context.Context, r CreateRequest) (*Product, error)
}

type Updater interface {
	Update(ctx context.Context, r UpdateRequest) (*Product, error)
}

type Deleter interface {
	Delete(ctx context.Context, id string) (*Product, error)
}
