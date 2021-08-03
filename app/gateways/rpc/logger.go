package rpc

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func loggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	l := log.With().Str("rpc", info.FullMethod).Logger()

	ctx = l.WithContext(ctx)

	return handler(ctx, req)
}
