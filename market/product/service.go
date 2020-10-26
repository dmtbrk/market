package product

import (
	"context"
	"fmt"
	"github.com/ortymid/market/market/auth"
)

type Service struct {
	Storage Storage
}

// List returns a list of products for the given request.
func (s *Service) Find(ctx context.Context, r FindRequest) ([]*Product, error) {
	ps, err := s.Storage.Find(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	return ps, nil
}

// Get returns a product for the given id. It returns product.ErrNotFound error if
// there is no product with such id.
func (s *Service) FindOne(ctx context.Context, id string) (*Product, error) {
	p, err := s.Storage.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}

	return p, nil
}

// Create creates a new product and returns it.
func (s *Service) Create(ctx context.Context, r CreateRequest) (p *Product, err error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("create product: %w", err)
	}

	r.Seller = user.ID
	p, err = s.Storage.Create(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return p, nil
}

// Update updates a product for the given id and returns it. It returns product.ErrNotFound
// error if there is no product with such id.
func (s *Service) Update(ctx context.Context, r UpdateRequest) (*Product, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("update product: %w", err)
	}

	p, err := s.Storage.FindOne(ctx, r.ID)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	if user.ID != p.Seller {
		err := auth.ErrPermission{Reason: "only own products allowed to update"}
		return nil, fmt.Errorf("update product: %w", err)
	}

	p, err = s.Storage.Update(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("update product: %w", err)
	}

	return p, nil
}

// Delete deletes a product for the given id. It returns product.ErrNotFound
// error if there is no product with such id.
func (s *Service) Delete(ctx context.Context, id string) (*Product, error) {
	user, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}
	if user == nil {
		err := auth.ErrPermission{Reason: "user not provided"}
		return nil, fmt.Errorf("delete product: %w", err)
	}

	p, err := s.Storage.FindOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}

	if user.ID != p.Seller {
		err := auth.ErrPermission{Reason: "only own products allowed to delete"}
		return nil, fmt.Errorf("delete product: %w", err)
	}

	p, err = s.Storage.Delete(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete product: %w", err)
	}

	return p, nil
}
