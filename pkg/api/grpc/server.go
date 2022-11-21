package grpc

import (
	"context"

	"github.com/opencars/grpc/pkg/core"
)

type vehicleHandler struct {
	core.UnimplementedVehicleServiceServer
	api *API
}

func (h *vehicleHandler) FindByNumber(ctx context.Context, r *core.NumberRequest) (*core.Result, error) {
	result, err := h.api.svc.FindByNumber(ctx, r.Number)
	if err != nil {
		return nil, err
	}

	return result.ToGRPC(), nil
}

func (h *vehicleHandler) FindByVIN(ctx context.Context, r *core.VINRequest) (*core.Result, error) {
	result, err := h.api.svc.FindByVIN(ctx, r.Vin)
	if err != nil {
		return nil, err
	}

	return result.ToGRPC(), nil
}
