package grpc

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type AuthFunc func(context.Context) (context.Context, error)

type AuthMiddleware struct {
	AuthFunc AuthFunc

	protectedMethods map[string]bool
}

func (m *AuthMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("JWT Interceptor: method %q", info.FullMethod)

		ctx, err := m.AuthFunc(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (m *AuthMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("JWT stream handling...")

		ctx := stream.Context()
		ctx, err := m.AuthFunc(ctx)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

// func (itr *AuthMiddleware) authorize(ctx context.Context, method string) error {
// 	if protected := itr.protectedMethods[method]; !protected {
// 		// Everyone can access this method.
// 		return nil
// 	}

// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
// 	}

// 	values := md["authorization"]
// 	if len(values) == 0 {
// 		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
// 	}

// 	tokenString := values[0]
// 	claims, err := jwt.Parse(tokenString, itr.Alg, itr.Secret)
// 	if err != nil {
// 		return status.Errorf(codes.Unauthenticated, "jwt is invalid: %v", err)
// 	}

// 	if len(claims.UserID) > 0 {
// 		// Access granted.
// 		ctx = context.WithValue(ctx, "UserID", claims.UserID)
// 		return nil
// 	}

// 	return status.Error(codes.PermissionDenied, "permission to access this method denied")
// }
