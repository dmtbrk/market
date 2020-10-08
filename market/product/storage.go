package product

import "context"

//go:generate mockgen -destination=../../mock/product_storage.go -package mock -mock_names=Storage=MockProductStorage . Storage

type Storage interface {
	Lister
	Getter
	Creator
	Updater
	Deleter
}

type Lister interface {
	List(ctx context.Context, r ListRequest) ([]*Product, error)
}

type Getter interface {
	Get(ctx context.Context, id string) (*Product, error)
}

type Creator interface {
	Create(ctx context.Context, r CreateRequest) (*Product, error)
}

type Updater interface {
	Update(ctx context.Context, r UpdateRequest) (*Product, error)
}

type Deleter interface {
	Delete(ctx context.Context, id string) (*Product, error)
}
