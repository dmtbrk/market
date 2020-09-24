package grpc

import (
	"context"
	"io"
	"log"

	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/market"
	"google.golang.org/grpc"
)

// Client implements market.Interface. It allows making calls to market
// gRPC server.
type Client struct {
	conn   *grpc.ClientConn
	client pb.MarketClient
}

// Connect must be called before any usage of Client. It connects to the
// market gRPC server by the provided address.
func (c *Client) Connect(serverAddr string) error {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.conn = conn

	client := pb.NewMarketClient(conn)
	c.client = client

	return nil
}

// Disconnect closes connection to the market gRPC server.
func (c *Client) Disconnect() error {
	return c.conn.Close()
}

func (c *Client) Products(ctx context.Context, offset int, limit int) ([]*market.Product, error) {
	req := &pb.ProductsRequest{
		Offset: int32(offset),
		Limit:  int32(limit),
	}

	stream, err := c.client.Products(ctx, req)
	if err != nil {
		return nil, err
	}

	var products []*market.Product
	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("receiving product:", err)
		}

		p := &market.Product{
			ID:     int(rep.Id),
			Name:   rep.Name,
			Price:  int(rep.Price),
			Seller: rep.Seller,
		}

		products = append(products, p)
	}

	return products, nil
}

func (c *Client) Product(ctx context.Context, id int) (*market.Product, error) {
	req := &pb.ProductRequest{
		Id: int32(id),
	}

	rep, err := c.client.Product(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &market.Product{
		ID:     int(rep.Id),
		Name:   rep.Name,
		Price:  int(rep.Price),
		Seller: rep.Seller,
	}
	return p, nil
}

func (c *Client) AddProduct(ctx context.Context, r market.AddProductRequest, userID string) (*market.Product, error) {
	req := &pb.AddProductRequest{
		Name:   r.Name,
		Price:  int32(r.Price),
		Seller: r.Seller,
		UserID: userID,
	}

	rep, err := c.client.AddProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &market.Product{
		ID:     int(rep.Id),
		Name:   rep.Name,
		Price:  int(rep.Price),
		Seller: rep.Seller,
	}
	return p, nil
}

func (c *Client) EditProduct(ctx context.Context, r market.EditProductRequest, userID string) (*market.Product, error) {
	req := &pb.EditProductRequest{
		Id:     int32(r.ID),
		Name:   r.Name,
		UserID: userID,
	}
	if r.Price != nil {
		price := int32(*r.Price)
		req.Price = &price
	}

	rep, err := c.client.EditProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	p := &market.Product{
		ID:     int(rep.Id),
		Name:   rep.Name,
		Price:  int(rep.Price),
		Seller: rep.Seller,
	}
	return p, nil
}

func (c *Client) DeleteProduct(ctx context.Context, id int, userID string) error {
	req := &pb.DeleteProductRequest{
		Id:     int32(id),
		UserID: userID,
	}

	_, err := c.client.DeleteProduct(ctx, req)

	return err
}
