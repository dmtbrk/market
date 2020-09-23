package grpc

import (
	"context"
	"io"
	"log"

	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/market"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MarketClient
}

func (c *Client) Connect(serverAddr string) error {
	conn, err := grpc.Dial(serverAddr)
	if err != nil {
		return err
	}
	c.conn = conn

	client := pb.NewMarketClient(conn)
	c.client = client

	return nil
}

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
	return nil, nil
}

func (c *Client) AddProduct(ctx context.Context, r market.AddProductRequest, userID string) (*market.Product, error) {
	return nil, nil
}

func (c *Client) EditProduct(ctx context.Context, r market.EditProductRequest, userID string) (*market.Product, error) {
	return nil, nil
}

func (c *Client) DeleteProduct(ctx context.Context, id int, userID string) error {
	return nil
}
