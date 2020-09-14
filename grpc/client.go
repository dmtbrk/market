package grpc

import (
	"context"
	"io"

	"github.com/ortymid/t3-grpc/grpc/pb"
	"github.com/ortymid/t3-grpc/market"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MarketClient
}

func NewMarket(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewMarketClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Products() ([]*market.Product, error) {
	stream, err := c.client.Products(context.TODO(), &pb.ProductsRequest{})
	if err != nil {
		return nil, err
	}

	ps := make([]*market.Product, 0)
	for {
		pr, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		ps = append(ps, &market.Product{
			ID:     int(pr.Id),
			Name:   pr.Name,
			Price:  int(pr.Price),
			Seller: pr.Seller,
		})
	}
	return ps, nil
}

func (c *Client) Product(id int) (*market.Product, error) {
	pr, err := c.client.Product(context.TODO(), &pb.ProductRequest{Id: int32(id)})
	if err != nil {
		return nil, err
	}
	return &market.Product{
		ID:     int(pr.Id),
		Name:   pr.Name,
		Price:  int(pr.Price),
		Seller: pr.Seller,
	}, nil
}

func (c *Client) AddProduct(p *market.Product, userID string) (*market.Product, error) {
	return &market.Product{}, nil
}
func (c *Client) ReplaceProduct(p *market.Product, userID string) (*market.Product, error) {
	return &market.Product{}, nil
}
func (c *Client) DeleteProduct(id int, userID string) error {
	return nil
}
