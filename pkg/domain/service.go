package domain

import (
	"context"

	"github.com/opencars/core/pkg/domain/model"
	"github.com/opencars/grpc/pkg/alpr"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/grpc/pkg/vin_decoding"
)

type RegistrationProvider interface {
	FindByVIN(ctx context.Context, vin string) ([]*registration.Record, error)
	FindByNumber(ctx context.Context, number string) ([]*registration.Record, error)
	FindByCode(ctx context.Context, code string) ([]*registration.Record, error)
}

type OperationProvider interface {
	FindByVIN(ctx context.Context, vin string) ([]*operation.Record, error)
	FindByNumber(ctx context.Context, number string) ([]*operation.Record, error)
}

type CoreService interface {
	FindByNumber(ctx context.Context, number string) (*model.Aggregate, error)
	FindByVIN(ctx context.Context, vin string) (*model.Aggregate, error)
}

type VinDecoder interface {
	Decode(context.Context, ...string) ([]*vin_decoding.DecodeResultItem, error)
}

type VehicleImageProvider interface {
	FindByNumber(context.Context, *alpr.NumberRequest) ([]*alpr.Recognition, error)
}
