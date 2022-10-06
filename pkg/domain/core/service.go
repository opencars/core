package core

import (
	"context"
	"time"

	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/operations/pkg/logger"

	"github.com/opencars/core/pkg/domain"
)

type Service struct {
	r  domain.RegistrationProvider
	o  domain.OperationProvider
	vd domain.VinDecoder
}

func NewService(r domain.RegistrationProvider, o domain.OperationProvider, vd domain.VinDecoder) (*Service, error) {
	return &Service{
		r:  r,
		o:  o,
		vd: vd,
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
				VIN: &core.Vin{
					Value: r.Vin,
				},
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

	if len(result.Vehicles) == 0 {
		for _, op := range operations {
			if _, ok := result.Vehicles[op.Vin]; !ok {
				v := domain.Vehicle{
					VIN: &core.Vin{
						Value: op.Vin,
					},
					Brand: op.Brand,
					Model: op.Model,
					Year:  op.Year,
				}

				result.Vehicles[op.Vin] = &v
			}

			result.Vehicles[op.Vin].Operations = append(result.Vehicles[op.Vin].Operations, op)
		}
	}

	vins := make([]string, 0, len(result.Vehicles))
	for vin := range result.Vehicles {
		vins = append(vins, vin)
	}

	decodedVins, err := s.vd.Decode(ctx, vins...)
	if err != nil {
		logger.Errorf("failed send decode command: %s", err)
	} else {
		for i, vinResult := range decodedVins {
			if vinResult.Error != nil {
				logger.Errorf("failed to parse vin code: %s", err)
				continue
			}

			vin := vins[i]
			vehicle, ok := result.Vehicles[vin]
			if !ok {
				logger.Errorf("failed to find : %s", err)
				continue
			}

			vehicle.VIN.DecodedVin = vinResult.DecodedVin
			vehicle.VIN.Vehicle = vinResult.Vehicle
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
				VIN:   &core.Vin{Value: r.Vin},
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
