package gql

import "github.com/ortymid/market/market/product"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProductService product.Interface
}
