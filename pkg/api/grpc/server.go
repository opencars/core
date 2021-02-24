package grpc

import (
	"context"

	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/core"
)

type vehicleHandler struct {
	core.UnimplementedVehicleServiceServer
	api *API
}

func (h *vehicleHandler) FindByNumber(ctx context.Context, r *core.NumberRequest) (*core.Result, error) {
	resp, err := h.api.svc.FindByNumber(ctx, r.Number)
	if err != nil {
		return nil, handleErr(err)
	}

	result := core.Result{
		Vehicles: make(map[string]*core.Vehicle),
	}

	for k, v := range resp.Vehicles {
		dto := core.Vehicle{
			Vin: v.VIN,
			FirstRegDate: &common.Date{
				Year:  int32(v.FirstRegDate.Year()),
				Month: int32(v.FirstRegDate.Month()),
				Day:   int32(v.FirstRegDate.Day()),
			},
			Brand: v.Brand,
			Model: v.Model,
			Year:  v.Year,
		}

		dto.Registrations = v.Registrations
		dto.Operations = v.Operations

		result.Vehicles[k] = &dto
	}

	return &result, nil
}

func (h *vehicleHandler) FindByVIN(ctx context.Context, r *core.VINRequest) (*core.Result, error) {
	resp, err := h.api.svc.FindByVIN(ctx, r.Vin)
	if err != nil {
		return nil, handleErr(err)
	}

	result := core.Result{
		Vehicles: make(map[string]*core.Vehicle),
	}

	for k, v := range resp.Vehicles {
		dto := core.Vehicle{
			Vin: v.VIN,
			FirstRegDate: &common.Date{
				Year:  int32(v.FirstRegDate.Year()),
				Month: int32(v.FirstRegDate.Month()),
				Day:   int32(v.FirstRegDate.Day()),
			},
			Brand: v.Brand,
			Model: v.Model,
			Year:  v.Year,
		}

		dto.Registrations = v.Registrations
		dto.Operations = v.Operations

		result.Vehicles[k] = &dto
	}

	return &result, nil
}
