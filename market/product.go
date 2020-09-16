package market

import (
	"context"
	"errors"
)

//go:generate mockgen -destination=./mock/product_service.go  -package=mock . ProductService

var ErrProductNotFound = errors.New("product not found")

// ProductService represents a product data backend.
type ProductService interface {
	Products(ctx context.Context, offset int, limit int) ([]*Product, error)
	Product(ctx context.Context, id int) (*Product, error)
	AddProduct(ctx context.Context, r AddProductRequest) (*Product, error)
	EditProduct(ctx context.Context, r EditProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id int) error

	ReplaceProduct(*Product) (*Product, error)
}

type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Seller string `json:"seller"`
}

type AddProductRequest struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Seller string `json:"seller"`
}

// EditProductRequest contains fields meant to be changed on the product with
// the specified id. Zero value represents no changes required for the field.
type EditProductRequest struct {
	ID     int    `json:"id"` // required
	Name   string `json:"name"`
	Price  *int   `json:"price"` // nil means do nothing with the field, zero means make it free
	Seller string `json:"seller"`
}
