package mem

import (
	"context"
	"sync"

	"github.com/ortymid/market/market"
)

type ProductService struct {
	mu       sync.RWMutex
	lastID   int
	products []*market.Product
}

func NewProductService() *ProductService {
	products := []*market.Product{
		{ID: 1, Name: "Banana", Price: 1500, Seller: "1234"},
		{ID: 2, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 3, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 4, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 5, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 6, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 7, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 8, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 9, Name: "Carrot", Price: 1400, Seller: "bunny"},
		{ID: 10, Name: "Carrot", Price: 1400, Seller: "bunny"},
	}
	return &ProductService{products: products, lastID: 2}
}

func (srv *ProductService) Products(ctx context.Context, offset int, limit int) ([]*market.Product, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	if offset >= len(srv.products) {
		return []*market.Product{}, nil
	}
	if offset < 0 {
		offset = 0
	}

	end := offset + limit
	if end > len(srv.products) {
		end = len(srv.products)
	}

	products := make([]*market.Product, end-offset)
	copy(products, srv.products[offset:end])

	return products, nil
}

func (srv *ProductService) Product(ctx context.Context, id int) (*market.Product, error) {
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	for _, p := range srv.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, market.ErrProductNotFound
}

func (srv *ProductService) AddProduct(ctx context.Context, r market.AddProductRequest) (*market.Product, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	srv.lastID++

	p := market.Product{
		ID:     srv.lastID,
		Name:   r.Name,
		Price:  r.Price,
		Seller: r.Seller,
	}

	srv.products = append(srv.products, &p)
	return &p, nil
}

func (srv *ProductService) EditProduct(ctx context.Context, r market.EditProductRequest) (*market.Product, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for i, p := range srv.products {
		if p.ID == r.ID {
			if len(r.Name) != 0 {
				p.Name = r.Name
			}
			if r.Price != nil {
				p.Price = *r.Price
			}
			if len(r.Seller) != 0 {
				p.Seller = r.Seller
			}
			srv.products[i] = p
			return p, nil
		}
	}
	return nil, market.ErrProductNotFound
}

func (srv *ProductService) ReplaceProduct(np *market.Product) (*market.Product, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for i, op := range srv.products {
		if op.ID == np.ID {
			srv.products[i] = np
			return np, nil
		}
	}
	return nil, market.ErrProductNotFound
}

func (srv *ProductService) DeleteProduct(ctx context.Context, id int) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	for i, p := range srv.products {
		if p.ID == id {
			if i == len(srv.products)-1 {
				srv.products[i] = nil
				srv.products = srv.products[:i]
			}
			copy(srv.products[i:], srv.products[i+1:])
			return nil
		}
	}
	return market.ErrProductNotFound
}
