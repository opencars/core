package domain

import (
	"context"

	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
)

type RegistrationProvider interface {
	FindByVIN(ctx context.Context, vin string) ([]*registration.Record, error)
	FindByNumber(ctx context.Context, number string) ([]*registration.Record, error)
	FindByCode(ctx context.Context, code string) ([]*registration.Record, error)
}

type OperationProvider interface {
	FindByNumber(ctx context.Context, number string) ([]*operation.Record, error)
}

type CoreService interface {
	FindByNumber(ctx context.Context, number string) (*Aggregate, error)
	FindByVIN(ctx context.Context, vin string) (*Aggregate, error)
}
