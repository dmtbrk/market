package gql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/ortymid/market/market/product"

	"github.com/ortymid/market/gql/gen"
	"github.com/ortymid/market/gql/model"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	req := product.CreateRequest{
		Name:  input.Name,
		Price: input.Price,
	}

	p, err := r.ProductService.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}, nil
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, input model.UpdateProduct) (*model.Product, error) {
	req := product.UpdateRequest{
		ID:    input.ID,
		Name:  input.Name,
		Price: input.Price,
	}

	p, err := r.ProductService.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}, nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (*model.Product, error) {
	p, err := r.ProductService.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}, nil
}

func (r *queryResolver) Products(ctx context.Context, offset int64, limit int64) ([]*model.Product, error) {
	req := product.FindRequest{
		Offset: offset,
		Limit:  limit,
	}

	products, err := r.ProductService.Find(ctx, req)
	if err != nil {
		return nil, err
	}

	ps := make([]*model.Product, len(products))
	for i, p := range products {
		ps[i] = &model.Product{
			ID:     p.ID,
			Name:   p.Name,
			Price:  p.Price,
			Seller: p.Seller,
		}
	}
	return ps, nil
}

func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	p, err := r.ProductService.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		ID:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}, nil
}

// Mutation returns gen.MutationResolver implementation.
func (r *Resolver) Mutation() gen.MutationResolver { return &mutationResolver{r} }

// Query returns gen.QueryResolver implementation.
func (r *Resolver) Query() gen.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
