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

func (s *Server) Find(r *pb.FindRequest, stream pb.ProductService_FindServer) error {
	ctx := context.TODO()

	var priceRange *product.PriceRange
	if r.PriceRange != nil {
		priceRange = &product.PriceRange{
			From: r.PriceRange.From,
			To:   r.PriceRange.To,
		}
	}

	fr := product.FindRequest{
		Offset:     r.Offset,
		Limit:      r.Limit,
		Name:       r.Name,
		PriceRange: priceRange,
	}

	ps, err := s.ProductService.Find(ctx, fr)
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

func (s *Server) FindOne(ctx context.Context, r *pb.FindOneRequest) (*pb.ProductReply, error) {
	p, err := s.ProductService.FindOne(ctx, r.Id)
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
