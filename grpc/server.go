package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ortymid/market/grpc/pb"
	"github.com/ortymid/market/jwt"
	"github.com/ortymid/market/market"
)

type Server struct {
	Market market.Interface

	JWTAlg    string
	JWTSecret interface{}

	grpcServer *grpc.Server
}

func (s *Server) Run(port int) error {
	var opts []grpc.ServerOption

	opts = append(opts, s.middlewares()...)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMarketServer(grpcServer, s)
	s.grpcServer = grpcServer

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	return grpcServer.Serve(ln)
}

func (s *Server) middlewares() []grpc.ServerOption {
	auth := &AuthMiddleware{
		AuthFunc: func(ctx context.Context) (context.Context, error) {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
			}

			values := md["authorization"]
			if len(values) == 0 {
				// No token means anonymous request.
				return ctx, nil
				// return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
			}

			tokenString := values[0]
			claims, err := jwt.Parse(tokenString, s.JWTAlg, s.JWTSecret)
			if err != nil {
				return ctx, status.Errorf(codes.Unauthenticated, "jwt is invalid: %v", err)
			}

			if len(claims.UserID) > 0 {
				// Access granted.
				ctx = context.WithValue(ctx, market.ContextKeyUserID, claims.UserID)
				return ctx, nil
			}

			return ctx, status.Error(codes.PermissionDenied, "permission to access this method denied")
		},
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(auth.Unary()),
		grpc.StreamInterceptor(auth.Stream()),
	}

	return opts
}

func (s *Server) Products(req *pb.ProductsRequest, stream pb.Market_ProductsServer) error {
	ctx := context.TODO()

	products, err := s.Market.Products(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return err
	}

	for _, p := range products {
		rep := &pb.ProductReply{
			Id:     int32(p.ID),
			Name:   p.Name,
			Price:  int32(p.Price),
			Seller: p.Seller,
		}
		if err := stream.Send(rep); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Product(ctx context.Context, req *pb.ProductRequest) (*pb.ProductReply, error) {
	p, err := s.Market.Product(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     int32(p.ID),
		Name:   p.Name,
		Price:  int32(p.Price),
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.ProductReply, error) {
	userID, _ := ctx.Value(market.ContextKeyUserID).(string)

	r := market.AddProductRequest{
		Name:   req.Name,
		Price:  int(req.Price),
		Seller: req.Seller,
	}

	p, err := s.Market.AddProduct(ctx, r, userID)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     int32(p.ID),
		Name:   p.Name,
		Price:  int32(p.Price),
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) EditProduct(ctx context.Context, req *pb.EditProductRequest) (*pb.ProductReply, error) {
	userID, _ := ctx.Value(market.ContextKeyUserID).(string)

	var price *int
	if req.Price != nil {
		p := int(*req.Price)
		price = &p
	}

	r := market.EditProductRequest{
		Name:  req.Name,
		Price: price,
	}

	p, err := s.Market.EditProduct(ctx, r, userID)
	if err != nil {
		return nil, err
	}

	rep := &pb.ProductReply{
		Id:     int32(p.ID),
		Name:   p.Name,
		Price:  int32(p.Price),
		Seller: p.Seller,
	}
	return rep, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	userID, _ := ctx.Value(market.ContextKeyUserID).(string)

	err := s.Market.DeleteProduct(ctx, int(req.Id), userID)

	return &pb.Empty{}, err
}
