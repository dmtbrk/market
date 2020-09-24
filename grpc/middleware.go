package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// ContextFunc is meant to be a pluggable function that somehow manipulates
// with the context.
type ContextFunc func(context.Context) (context.Context, error)

// ContextMiddleware holds a function that is called on the request context.
type ContextMiddleware struct {
	ContextFunc ContextFunc
}

// Unary returns an unary interceptor.
func (m *ContextMiddleware) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctx, err := m.ContextFunc(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// Stream returns a stream interceptor.
func (m *ContextMiddleware) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()

		ctx, err := m.ContextFunc(ctx)
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
