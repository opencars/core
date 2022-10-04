package vin_decoding

import (
	"context"

	"google.golang.org/grpc"

	"github.com/opencars/grpc/pkg/vin_decoding"
)

type Service struct {
	c vin_decoding.ServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: vin_decoding.NewServiceClient(conn),
	}, nil
}

func (s *Service) Decode(ctx context.Context, vins ...string) ([]*vin_decoding.DecodeResultItem, error) {
	req := vin_decoding.DecodeRequest{
		Vins: vins,
	}

	resp, err := s.c.Decode(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}
