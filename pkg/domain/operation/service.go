package operation

import (
	"context"

	"google.golang.org/grpc"

	"github.com/opencars/grpc/pkg/operation"
)

type Service struct {
	c operation.ServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: operation.NewServiceClient(conn),
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) ([]*operation.Record, error) {
	resp, err := s.c.FindByNumber(ctx, &operation.NumberRequest{Number: number})
	if err != nil {
		return nil, err
	}

	return resp.Operations, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) ([]*operation.Record, error) {
	resp, err := s.c.FindByVIN(ctx, &operation.VINRequest{Vin: vin})
	if err != nil {
		return nil, err
	}

	return resp.Operations, nil
}
