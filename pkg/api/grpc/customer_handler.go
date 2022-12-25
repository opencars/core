package grpc

import (
	"context"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/core/pkg/domain/query"
	"github.com/opencars/grpc/pkg/core/customer"
)

type customerHandler struct {
	customer.UnimplementedVehicleServiceServer
	svc domain.CustomerService
}

func newCustomerHandler(svc domain.CustomerService) *customerHandler {
	return &customerHandler{
		svc: svc,
	}
}

func (h *customerHandler) FindByNumber(ctx context.Context, r *customer.FinByNumberRequest) (*customer.FindByNumberResponse, error) {
	q := query.ListByNumber{
		UserID:  UserIDFromContext(ctx),
		TokenID: TokenIDFromContext(ctx),
		Number:  r.Number,
	}

	result, err := h.svc.FindByNumber(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	vehicles := make([]*customer.Vehicle, 0, len(result.Vehicles))
	for _, v := range result.Vehicles {
		vehicles = append(vehicles, v.ToCustomerGRPC())
	}

	return &customer.FindByNumberResponse{
		Vehicles: vehicles,
	}, nil
}

func (h *customerHandler) FindByVIN(ctx context.Context, r *customer.FindByVinRequest) (*customer.FindByVinResponse, error) {
	q := query.ListByVIN{
		UserID:  UserIDFromContext(ctx),
		TokenID: TokenIDFromContext(ctx),
		VIN:     r.Vin,
	}

	result, err := h.svc.FindByVIN(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	return &customer.FindByVinResponse{
		Vehicle: result.ToCustomerGRPC(),
	}, nil
}
