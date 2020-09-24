package market

import (
	"context"
	"errors"
	"fmt"
)

//go:generate mockgen -destination=./mock/market.go -package=mock -mock_names=Interface=MockMarket . Interface

// ErrPermission is an error returned when a user does not have rights
// to do Market methods.
type ErrPermission struct {
	Reason error
}

func (err *ErrPermission) Error() string {
	return fmt.Sprintf("permission denied: %s", err.Reason)
}

func (err *ErrPermission) Is(target error) bool {
	_, ok := target.(*ErrPermission)
	return ok
}

func (err *ErrPermission) Unwrap() error {
	return err.Reason
}

// Interface may be used by protocol layers for RPC or mocking.
type Interface interface {
	Products(ctx context.Context, offset int, limit int) ([]*Product, error)
	Product(ctx context.Context, id int) (*Product, error)
	AddProduct(ctx context.Context, r AddProductRequest, userID string) (*Product, error)
	EditProduct(ctx context.Context, r EditProductRequest, userID string) (*Product, error)
	DeleteProduct(ctx context.Context, id int, userID string) error
}

// Market composes business logic from different services.
type Market struct {
	UserService    UserService
	ProductService ProductService
}

// Products returns all products on the market.
func (m *Market) Products(ctx context.Context, offset int, limit int) ([]*Product, error) {
	ps, err := m.ProductService.Products(ctx, offset, limit)
	if err != nil {
		err = fmt.Errorf("products: %w", err)
		return nil, err
	}
	return ps, nil
}

// Product finds the product by its ID.
func (m *Market) Product(ctx context.Context, id int) (*Product, error) {
	p, err := m.ProductService.Product(ctx, id)
	if err != nil {
		err = fmt.Errorf("product: %w", err)
		return nil, err
	}
	return p, nil
}

func (m *Market) AddProduct(ctx context.Context, r AddProductRequest, userID string) (*Product, error) {
	// Check the user for permission. Only the existence of the user counts yet.
	_, err := m.UserService.User(userID)
	if errors.Is(err, &ErrUserNotFound{}) {
		err = fmt.Errorf("add product: %w", &ErrPermission{Reason: err})
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("add product: %w", err)
		return nil, err
	}

	// After the user check, add the product.
	p, err := m.ProductService.AddProduct(ctx, r)
	if err != nil {
		err = fmt.Errorf("add product: %w", err)
		return nil, err
	}
	return p, nil
}

func (m *Market) EditProduct(ctx context.Context, r EditProductRequest, userID string) (*Product, error) {
	// Check the user for permission. Only the existence of the user counts yet.
	_, err := m.UserService.User(userID)
	if errors.Is(err, &ErrUserNotFound{}) {
		err = fmt.Errorf("edit product: %w", &ErrPermission{Reason: err})
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("edit product: %w", err)
		return nil, err
	}

	// After the user check, add the product.
	p, err := m.ProductService.EditProduct(ctx, r)
	if err != nil {
		err = fmt.Errorf("edit product: %w", err)
		return nil, err
	}
	return p, nil
}

// DeleteProduct deletes the product from the market by its ID
// checking the permission to do it by user ID.
func (m *Market) DeleteProduct(ctx context.Context, id int, userID string) error {
	// Obtain the product.
	product, err := m.ProductService.Product(ctx, id)
	if err != nil {
		err = fmt.Errorf("delete product: %w", err)
		return err
	}

	// Check the user for permission. The user existence does not count here.
	if product.Seller != userID {
		err = fmt.Errorf("delete product: %w", &ErrPermission{Reason: err})
		return err
	}

	// Perform deletion.
	err = m.ProductService.DeleteProduct(ctx, id)
	if err != nil {
		err = fmt.Errorf("delete product: %w", err)
		return err
	}
	return nil
}
