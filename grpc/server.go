package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ortymid/t3-grpc/grpc/pb"
	"github.com/ortymid/t3-grpc/market"
	"google.golang.org/grpc"
)

type Server struct {
	Market market.Interface
	server *grpc.Server
}

func (s *Server) Products(req *pb.ProductsRequest, stream pb.Market_ProductsServer) error {
	products, err := s.Market.Products()
	if err != nil {
		return err
	}

	for _, product := range products {
		pr := &pb.ProductReply{
			Id:     int32(product.ID),
			Name:   product.Name,
			Price:  int32(product.Price),
			Seller: product.Seller,
		}
		stream.Send(pr)
	}
	return nil
}

func (s *Server) Product(ctx context.Context, req *pb.ProductRequest) (*pb.ProductReply, error) {
	product, err := s.Market.Product(int(req.Id))
	if err != nil {
		return nil, err
	}

	pr := &pb.ProductReply{
		Id:     int32(product.ID),
		Name:   product.Name,
		Price:  int32(product.Price),
		Seller: product.Seller,
	}
	return pr, nil
}

func (s *Server) Run(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMarketServer(grpcServer, s)

	grpcServer.Serve(ln)

	s.server = grpcServer
	return nil
}

func (s *Server) Stop() {
	s.server.Stop()
}
