package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/opencars/core/pkg/domain/model"
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
