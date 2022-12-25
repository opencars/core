package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/opencars/seedwork/logger"
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

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !strings.Contains(info.FullMethod, "/core.customer.VehicleService") {
		return handler(ctx, req)
	}

	md, _ := metadata.FromIncomingContext(ctx)

	logger.Infof("md: %+v", md)

	userID := md.Get(HeaderUserID)
	if len(userID) == 0 {
		return nil, errors.New("auth: expected user id")
	}

	tokenID := md.Get(HeaderTokenID)
	if len(tokenID) == 0 {
		return nil, errors.New("auth: expected token id")
	}

	tokenName := md.Get(HeaderTokenName)
	if len(tokenName) == 0 {
		return nil, errors.New("auth: expected token name")
	}

	ctx = WithUserID(ctx, userID[0])
	ctx = WithTokenID(ctx, tokenID[0])
	ctx = WithTokenName(ctx, tokenName[0])

	return handler(ctx, req)
}
