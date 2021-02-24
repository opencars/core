package registration

import (
	"context"

	"github.com/opencars/grpc/pkg/registration"
	"google.golang.org/grpc"
)

type Service struct {
	c registration.ServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: registration.NewServiceClient(conn),
	}, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) ([]*registration.Record, error) {
	resp, err := s.c.FindByVIN(ctx, &registration.VINRequest{Vin: vin})
	if err != nil {
		return nil, err
	}

	return resp.Registrations, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) ([]*registration.Record, error) {
	resp, err := s.c.FindByNumber(ctx, &registration.NumberRequest{Number: number})
	if err != nil {
		return nil, err
	}

	return resp.Registrations, nil
}

func (s *Service) FindByCode(ctx context.Context, code string) ([]*registration.Record, error) {
	resp, err := s.c.FindByCode(ctx, &registration.CodeRequest{Code: code})
	if err != nil {
		return nil, err
	}

	return resp.Registrations, nil
}
