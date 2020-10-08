package grpc

import (
	"context"
	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/market/product"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	AuthService    AuthService
	ProductService product.Interface
}

func (s *Server) List(r *pb.ListRequest, stream pb.ProductService_ListServer) error {
	ctx := context.TODO()
	lr := product.ListRequest{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	ps, err := s.ProductService.List(ctx, lr)
	if err != nil {
		return nil
	}

	for _, p := range ps {
		rep := &pb.ProductReply{
			Id:     p.ID,
			Name:   p.Name,
			Price:  p.Price,
			Seller: p.Seller,
		}
		if err := stream.Send(rep); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Get(ctx context.Context, r *pb.GetRequest) (*pb.ProductReply, error) {
	p, err := s.ProductService.Get(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) Create(ctx context.Context, r *pb.CreateRequest) (*pb.ProductReply, error) {
	cr := product.CreateRequest{
		Name:  r.Name,
		Price: r.Price,
	}

	p, err := s.ProductService.Create(ctx, cr)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.ProductReply, error) {
	ur := product.UpdateRequest{
		ID:    r.Id,
		Name:  r.Name,
		Price: r.Price,
	}

	p, err := s.ProductService.Update(ctx, ur)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.ProductReply, error) {
	p, err := s.ProductService.Delete(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     p.ID,
		Name:   p.Name,
		Price:  p.Price,
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) Run(addr string) error {
	auth := AuthInterceptor{AuthService: s.AuthService}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
		grpc.StreamInterceptor(auth.StreamServerInterceptor()),
	)
	pb.RegisterProductServiceServer(grpcServer, s)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return grpcServer.Serve(ln)
}
