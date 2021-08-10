package rpc

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func loggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Strips the rpc method name from FullMethod. Returns the last part of info.FullMethod after '/'.
	l := log.With().Str("handler", info.FullMethod[strings.LastIndex(info.FullMethod, "/")+1:]).Logger()

	ctx = l.WithContext(ctx)

	return handler(ctx, req)
}
