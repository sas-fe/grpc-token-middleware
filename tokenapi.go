package tokenapi

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// TokenAPI interface implements methods for checking token validity and incrementing usage.
type TokenAPI interface {
	CheckValidity(ctx context.Context) (bool, error)
	IncrementUsage(ctx context.Context) (bool, error)
	AsyncIncrementUsage(ctx context.Context)
}

// UnaryServerInterceptor returns a new unary server interceptor that
// checks token validity per-request and increments usage if valid.
func UnaryServerInterceptor(tokenAPI TokenAPI, asyncInc bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		allowed, err := tokenAPI.CheckValidity(ctx)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, grpc.Errorf(codes.ResourceExhausted, "API limit exceeded")
		}

		resp, e := handler(ctx, req)
		if e != nil {
			return nil, e
		}

		if asyncInc {
			go tokenAPI.AsyncIncrementUsage(ctx)
		} else {
			_, err = tokenAPI.IncrementUsage(ctx)
			if err != nil {
				return nil, err
			}
		}

		return resp, e
	}
}

// StreamServerInterceptor returns a new stream server interceptor that
// checks token validity per-request and increments usage if valid.
func StreamServerInterceptor(tokenAPI TokenAPI, asyncInc bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		allowed, err := tokenAPI.CheckValidity(ss.Context())
		if err != nil {
			return err
		}
		if !allowed {
			return grpc.Errorf(codes.ResourceExhausted, "API limit exceeded")
		}

		e := handler(srv, ss)
		if e != nil {
			return e
		}

		if asyncInc {
			go tokenAPI.AsyncIncrementUsage(ss.Context())
		} else {
			_, err = tokenAPI.IncrementUsage(ss.Context())
			if err != nil {
				return err
			}
		}

		return e
	}
}
