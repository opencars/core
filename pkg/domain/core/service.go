package core

import (
	"context"
	"time"

	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/seedwork/logger"

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
	logger.Debugf("find all registratiton with given number")

	// Find all registratiton with given number.
	registrations, err := s.r.FindByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	logger.Debugf("find all operations with given number")

	// Find all operations with given number.
	operations, err := s.o.FindByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	logger.Debugf("detect all unique vehicles from given operations and registrations.")

	// Detect all unique vehicles from given operations and registrations.
	vehicles, err := s.detectVehicles(ctx, operations, registrations)
	if err != nil {
		return nil, err
	}

	logger.Debugf("collect operations and registrations for each vehicle vin.")

	// For each unique vehicle we loop throught existing operations
	// and try to find operations and registrations by vehicles vin.
	for k, v := range vehicles.Vehicles {
		if !v.HasVIN() {
			logger.Debugf("vehicle %s does not have vin-code", k)
			continue
		}

		// Find all operations/registrations with given vin-code.
		logger.Debugf("vehicle %s, registrations", v.VIN.GetValue())

		registrations, err := s.r.FindByVIN(ctx, v.VIN.GetValue())
		if err != nil {
			return nil, err
		}

		logger.Debugf("vehicle %s, operations", v.VIN.GetValue())

		operations, err := s.o.FindByVIN(ctx, v.VIN.GetValue())
		if err != nil {
			return nil, err
		}

		v.AppendOperations(operations...)
		v.AppendRegistrations(registrations...)
	}

	logger.Debugf("map all vins")

	vins := vehicles.VINs()

	logger.Debugf("decode each unique vin")

	// Decode each unique vin.
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
			vehicle, ok := vehicles.Vehicles[vin]
			if !ok {
				logger.Errorf("failed to find : %s", err)
				continue
			}

			vehicle.VIN.DecodedVin = vinResult.DecodedVin
			vehicle.VIN.Vehicle = vinResult.Vehicle
		}
	}

	return vehicles, nil
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

func (s *Service) detectVehicles(ctx context.Context, operations []*operation.Record, registrations []*registration.Record) (*domain.Aggregate, error) {
	logger.Debugf("detectVehicles")

	result := domain.Aggregate{
		Vehicles: make(map[string]*domain.Vehicle),
	}

	for _, r := range registrations {
		logger.Debugf("detectVehicles: registration: %#v", r.String())

		hash := domain.Hash(r)

		if _, ok := result.Vehicles[hash]; !ok {
			v := domain.NewVehicle(r.Vin, r.Brand, r.Model, r.Year)

			// Convert date of the first vehicle registration.
			if r.FirstRegDate != nil {
				firstRegDate := time.Date(
					int(r.FirstRegDate.Year),
					time.Month(r.FirstRegDate.Month),
					int(r.FirstRegDate.Day),
					0, 0, 0, 0,
					time.UTC,
				)

				v.SetFirstRegDate(firstRegDate)
			}

			result.Vehicles[hash] = &v
		}

		result.Vehicles[hash].AppendRegistrations(r)
	}

	for _, op := range operations {
		logger.Debugf("detectVehicles: operation: %#v", op.String())

		hash := domain.Hash(op)

		if _, ok := result.Vehicles[hash]; !ok {
			v := domain.NewVehicle(hash, op.Brand, op.Model, op.Year)
			result.Vehicles[hash] = &v
		}

		result.Vehicles[hash].AppendOperations(op)
	}

	return &result, nil
}
