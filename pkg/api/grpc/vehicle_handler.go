package grpc

import (
	"context"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/grpc/pkg/core"
)

type vehicleHandler struct {
	core.UnimplementedVehicleServiceServer
	svc domain.CoreService
}

func newVehicleHandler(svc domain.CoreService) *vehicleHandler {
	return &vehicleHandler{
		svc: svc,
	}
}

func (h *vehicleHandler) FindByNumber(ctx context.Context, r *core.NumberRequest) (*core.Result, error) {
	result, err := h.svc.FindByNumber(ctx, r.Number)
	if err != nil {
		return nil, handleErr(err)
	}

	return result.ToGRPC(), nil
}

func (h *vehicleHandler) FindByVIN(ctx context.Context, r *core.VINRequest) (*core.Result, error) {
	result, err := h.svc.FindByVIN(ctx, r.Vin)
	if err != nil {
		return nil, handleErr(err)
	}

	return result.ToGRPC(), nil
}
