package core

import (
	"context"
	"time"

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

			for _, v := range vehicles.Vehicles {
				if v.VIN.GetValue() == vins[i] {
					v.VIN.DecodedVin = vinResult.DecodedVin
					v.VIN.Vehicle = vinResult.Vehicle
				}
			}
		}
	}

	return vehicles, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (*domain.Aggregate, error) {
	logger.Debugf("find all registratiton with given vin")

	// Find all registratiton with given vin.
	registrations, err := s.r.FindByVIN(ctx, vin)
	if err != nil {
		return nil, err
	}

	logger.Debugf("find all operations with given vin")

	// Find all operations with given vin.
	operations, err := s.o.FindByVIN(ctx, vin)
	if err != nil {
		return nil, err
	}

	logger.Debugf("detect all unique vehicles from given operations and registrations.")

	// Detect all unique vehicles from given operations and registrations.
	vehicles, err := s.detectVehicles(ctx, operations, registrations)
	if err != nil {
		return nil, err
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

			for _, v := range vehicles.Vehicles {
				if v.VIN.GetValue() == vins[i] {
					v.VIN.DecodedVin = vinResult.DecodedVin
					v.VIN.Vehicle = vinResult.Vehicle
				}
			}
		}
	}

	return vehicles, nil
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
