package grpc

import (
	"context"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/grpc/pkg/core"
)

type customerHandler struct {
	core.UnimplementedCustomerServiceServer
	svc domain.CustomerService
}

func newCustomerHandler(svc domain.CustomerService) *customerHandler {
	return &customerHandler{
		svc: svc,
	}
}

func (h *customerHandler) FindByNumber(ctx context.Context, r *core.NumberRequestByCustomer) (*core.ResultForCustomer, error) {
	result, err := h.svc.FindByNumber(ctx, r.Number)
	if err != nil {
		return nil, err
	}

	return result.ToExternalGRPC(), nil
}

func (h *customerHandler) FindByVIN(ctx context.Context, r *core.VINRequestByCustomer) (*core.ResultForCustomer, error) {
	result, err := h.svc.FindByVIN(ctx, r.Vin)
	if err != nil {
		return nil, err
	}

	return result.ToExternalGRPC(), nil
}
