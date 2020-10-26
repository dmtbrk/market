package grpc

import (
	"context"
	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/market/product"
	"google.golang.org/grpc"
	"io"
	"log"
)

// ProductService implements product.Interface. It allows making calls to the market
// gRPC server.
type ProductService struct {
	AuthService AuthService

	client pb.ProductServiceClient
}

func NewProductService(auth AuthService) *ProductService {
	return &ProductService{AuthService: auth}
}

// Connect must be called before any usage of Client. It connects to the
// market gRPC server at the provided address.
func (s *ProductService) Connect(ctx context.Context, addr string) error {
	auth := AuthInterceptor{AuthService: s.AuthService}

	conn, err := grpc.DialContext(
		ctx, addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(auth.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(auth.StreamClientInterceptor()),
	)
	if err != nil {
		return err
	}

	s.client = pb.NewProductServiceClient(conn)

	return nil
}

func (s *ProductService) Find(ctx context.Context, r product.FindRequest) ([]*product.Product, error) {
	req := &pb.ListRequest{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	stream, err := s.client.List(ctx, req)
	if err != nil {
		return nil, err
	}

	products := make([]*product.Product, 0, r.Limit)
	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("receiving product:", err)
		}

		p := &product.Product{
			ID:     rep.Id,
			Name:   rep.Name,
			Price:  rep.Price,
			Seller: rep.Seller,
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *ProductService) FindOne(ctx context.Context, id string) (*product.Product, error) {
	req := &pb.GetRequest{
		Id: id,
	}

	rep, err := s.client.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &product.Product{
		ID:     rep.Id,
		Name:   rep.Name,
		Price:  rep.Price,
		Seller: rep.Seller,
	}
	return p, nil
}

func (s *ProductService) Create(ctx context.Context, r product.CreateRequest) (*product.Product, error) {
	req := &pb.CreateRequest{
		Name:  r.Name,
		Price: r.Price,
	}

	rep, err := s.client.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &product.Product{
		ID:     rep.Id,
		Name:   rep.Name,
		Price:  rep.Price,
		Seller: rep.Seller,
	}
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, r product.UpdateRequest) (*product.Product, error) {
	req := &pb.UpdateRequest{
		Id:   r.ID,
		Name: r.Name,
	}
	if r.Price != nil {
		price := *r.Price
		req.Price = &price
	}

	rep, err := s.client.Update(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &product.Product{
		ID:     rep.Id,
		Name:   rep.Name,
		Price:  rep.Price,
		Seller: rep.Seller,
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) (*product.Product, error) {
	req := &pb.DeleteRequest{
		Id: id,
	}

	rep, err := s.client.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &product.Product{
		ID:     rep.Id,
		Name:   rep.Name,
		Price:  rep.Price,
		Seller: rep.Seller,
	}
	return p, nil
}
