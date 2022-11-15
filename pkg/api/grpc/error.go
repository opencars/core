package grpc

import (
	"github.com/opencars/core/pkg/domain/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound = status.Error(codes.NotFound, "record.not_found")
)

func handleErr(err error) error {
	switch err {
	case model.ErrNotFound:
		return ErrNotFound
	default:
		return err
	}
}
