package alpr

import (
	"context"

	"google.golang.org/grpc"

	"github.com/opencars/grpc/pkg/alpr"
)

type Service struct {
	c alpr.ServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: alpr.NewServiceClient(conn),
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, r *alpr.NumberRequest) ([]*alpr.Recognition, error) {
	req := alpr.NumberRequest{
		Number: r.Number,
	}

	resp, err := s.c.FindByNumber(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.Recognitions, nil
}
