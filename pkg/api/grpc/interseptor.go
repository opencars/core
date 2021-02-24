package grpc

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc"

	"github.com/opencars/core/pkg/logger"
)

// RequestLoggingInterceptor write request body to logs.
func RequestLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log := logger.WithFields(logger.Fields{
		"method": info.FullMethod,
	})

	reqBody, err := json.Marshal(req)
	if err != nil {
		log.Errorf("failed to unmarshal request: %s", err)
		return nil, err
	}

	log.Debugf("start handling new request: %s", reqBody)

	return handler(ctx, req)
}
