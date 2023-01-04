package wanted

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/opencars/grpc/pkg/wanted"
	"github.com/opencars/seedwork/logger"
)

type Service struct {
	c wanted.ServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Service{
		c: wanted.NewServiceClient(conn),
	}, nil
}

func (s *Service) Find(ctx context.Context, vins, numbers []string) ([]*wanted.Vehicle, error) {
	req := wanted.FindRequest{
		Numbers: numbers,
		Vins:    vins,
	}

	resp, err := s.c.Find(ctx, &req)
	if err != nil {
		return nil, err
	}

	logger.Debugf("wanted resp: %+v", resp.Vehicles)

	return resp.Vehicles, nil
}
