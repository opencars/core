package grpc

import (
	"context"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/core/pkg/domain/query"
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
	q := query.ListByNumber{
		UserID:  UserIDFromContext(ctx),
		TokenID: TokenIDFromContext(ctx),
		Number:  r.Number,
	}

	result, err := h.svc.FindByNumber(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	return result.ToExternalGRPC(), nil
}

func (h *customerHandler) FindByVIN(ctx context.Context, r *core.VINRequestByCustomer) (*core.ResultForCustomer, error) {
	q := query.ListByVIN{
		UserID:  UserIDFromContext(ctx),
		TokenID: TokenIDFromContext(ctx),
		VIN:     r.Vin,
	}

	result, err := h.svc.FindByVIN(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	return result.ToExternalGRPC(), nil
}
