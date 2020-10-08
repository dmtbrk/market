package grpc

import (
	"context"
	"github.com/ortymid/market/market/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	AuthService AuthService
}

func (s *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		u, err := s.AuthService.Authorize(ctx, md)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx = auth.NewContextWithUser(ctx, u)
		return handler(ctx, req)
	}
}

func (s *AuthInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		u, err := s.AuthService.Authorize(ctx, md)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx = auth.NewContextWithUser(ctx, u)
		return handler(srv, ss)
	}
}

func (s *AuthInterceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		u, err := auth.UserFromContext(ctx)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		md, err := s.AuthService.MetadataWithAuthorization(ctx, u)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (s *AuthInterceptor) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		u, err := auth.UserFromContext(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		md, err := s.AuthService.MetadataWithAuthorization(ctx, u)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
