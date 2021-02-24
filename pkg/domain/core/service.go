package core

import (
	"context"
	"time"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/grpc/pkg/operation"
)

type Service struct {
	r domain.RegistrationProvider
	o domain.OperationProvider
}

func NewService(r domain.RegistrationProvider, o domain.OperationProvider) (*Service, error) {
	return &Service{
		r: r,
		o: o,
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) (*domain.Aggregate, error) {
	registrations, err := s.r.FindByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	operations, err := s.o.FindByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	result := domain.Aggregate{
		Vehicles: make(map[string]*domain.Vehicle),
	}

	for _, r := range registrations {
		if _, ok := result.Vehicles[r.Vin]; !ok {
			v := domain.Vehicle{
				VIN:   r.Vin,
				Brand: r.Brand,
				Model: r.Model,
				Year:  r.Year,
			}

			// Convert date of the first vehicle registration.
			if r.FirstRegDate != nil {
				firstRegDate := time.Date(
					int(r.FirstRegDate.Year),
					time.Month(r.FirstRegDate.Month),
					int(r.FirstRegDate.Day),
					0, 0, 0, 0,
					time.UTC,
				)
				v.FirstRegDate = &firstRegDate
			}

			result.Vehicles[r.Vin] = &v
		}

		result.Vehicles[r.Vin].Registrations = append(result.Vehicles[r.Vin].Registrations, r)
	}

	for _, v := range result.Vehicles {
		for _, o := range operations {
			if o.Model == v.Model && o.Year == v.Year {
				v.Operations = append(v.Operations, o)
			}
		}
	}

	return &result, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (*domain.Aggregate, error) {
	registrations, err := s.r.FindByVIN(ctx, vin)
	if err != nil {
		return nil, err
	}

	result := domain.Aggregate{
		Vehicles: make(map[string]*domain.Vehicle),
	}

	for _, r := range registrations {
		if _, ok := result.Vehicles[r.Vin]; !ok {
			v := domain.Vehicle{
				VIN:   r.Vin,
				Brand: r.Brand,
				Model: r.Model,
				Year:  r.Year,
			}

			// Convert date of the first vehicle registration.
			if r.FirstRegDate != nil {
				firstRegDate := time.Date(
					int(r.FirstRegDate.Year),
					time.Month(r.FirstRegDate.Month),
					int(r.FirstRegDate.Day),
					0, 0, 0, 0,
					time.UTC,
				)
				v.FirstRegDate = &firstRegDate
			}

			result.Vehicles[r.Vin] = &v
		}

		result.Vehicles[r.Vin].Registrations = append(result.Vehicles[r.Vin].Registrations, r)
	}

	for _, v := range result.Vehicles {
		numbers := make(map[string]struct{})

		for _, r := range registrations {
			numbers[r.Number] = struct{}{}
		}

		allOperations := make([]*operation.Record, 0)
		for number := range numbers {
			operations, err := s.o.FindByNumber(ctx, number)
			if err != nil {
				return nil, err
			}

			for _, o := range operations {
				if v.Brand == o.Brand && v.Year == o.Year {
					allOperations = append(allOperations, o)
				}
			}
		}

		v.Operations = append(v.Operations, allOperations...)
	}

	return &result, nil
}
