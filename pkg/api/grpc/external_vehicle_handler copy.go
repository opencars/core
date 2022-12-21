package grpc

import (
	"context"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/grpc/pkg/core"
)

type externalVehicleHandler struct {
	core.UnimplementedExternalVehicleServiceServer
	svc domain.ExternalService
}

func newExternalVehicleHandler(svc domain.ExternalService) *externalVehicleHandler {
	return &externalVehicleHandler{
		svc: svc,
	}
}

func (h *externalVehicleHandler) FindByNumber(ctx context.Context, r *core.ExternalNumberRequest) (*core.ExternalResult, error) {
	result, err := h.svc.FindByNumber(ctx, r.Number)
	if err != nil {
		return nil, err
	}

	return result.ToExternalGRPC(), nil
}

func (h *externalVehicleHandler) FindByVIN(ctx context.Context, r *core.ExternalVINRequest) (*core.ExternalResult, error) {
	result, err := h.svc.FindByVIN(ctx, r.Vin)
	if err != nil {
		return nil, err
	}

	return result.ToExternalGRPC(), nil
}
