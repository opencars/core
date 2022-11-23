package core

import (
	"context"
	"time"

	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/seedwork/logger"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/core/pkg/domain/model"
)

type Service struct {
	r  domain.RegistrationProvider
	o  domain.OperationProvider
	vd domain.VinDecoder
	as domain.AdvertisementService
}

func NewService(r domain.RegistrationProvider, o domain.OperationProvider, vd domain.VinDecoder, as domain.AdvertisementService) (*Service, error) {
	return &Service{
		r:  r,
		o:  o,
		vd: vd,
		as: as,
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) (*model.Aggregate, error) {
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
	for k, v := range vehicles {
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

	vins := model.GetVINs(vehicles)

	logger.Debugf("decode each unique vin")

	// Search information about vehicle adds.
	advertisements, err := s.as.FindByVINs(ctx, vins, []string{})
	if err != nil {
		logger.Errorf("failed to get adds: %s", err)
	} else {
		for _, add := range advertisements {
			logger.Errorf("addv: %+v", add)
			for _, v := range vehicles {
				if add.VinPage != v.VIN.Value && add.VinOpencars != v.VIN.Value {
					logger.Errorf("unexpected vin: %s", add.VinPage)
					continue
				}

				v.AppendAdvertisements(add)
			}
		}
	}

	// Decode each unique vin.
	decodedVins, err := s.vd.Decode(ctx, vins...)
	if err != nil {
		logger.Errorf("failed send decode command: %s", err)
	} else {
		for i, vinResult := range decodedVins {
			if vinResult.Error != nil {
				logger.Errorf("failed to parse vin-code %s: %s", vins[i], vinResult.Error)
				continue
			}

			for _, v := range vehicles {
				if v.VIN.GetValue() == vins[i] {
					v.VIN.DecodedVin = vinResult.DecodedVin
					v.VIN.Vehicle = vinResult.Vehicle
				}
			}
		}
	}

	return model.NewAggregateWithNumber(number, vehicles), nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (*model.Aggregate, error) {
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

	vins := model.GetVINs(vehicles)

	logger.Debugf("decode each unique vin")

	// Search information about vehicle adds.
	advertisements, err := s.as.FindByVINs(ctx, vins, []string{})
	if err != nil {
		logger.Errorf("failed to get adds: %s", err)
	} else {
		for _, add := range advertisements {
			logger.Errorf("addv: %+v", add)
			for _, v := range vehicles {
				if add.VinPage != v.VIN.Value && add.VinOpencars != v.VIN.Value {
					logger.Errorf("unexpected vin: %s", add.VinPage)
					continue
				}

				v.AppendAdvertisements(add)
			}
		}
	}

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

			for _, v := range vehicles {
				if v.VIN.GetValue() == vins[i] {
					v.VIN.DecodedVin = vinResult.DecodedVin
					v.VIN.Vehicle = vinResult.Vehicle
				}
			}
		}
	}

	return model.NewAggregate(vehicles), nil
}

func (s *Service) detectVehicles(ctx context.Context, operations []*operation.Record, registrations []*registration.Record) (map[string]*model.Vehicle, error) {
	vehicles := make(map[string]*model.Vehicle)

	for _, r := range registrations {
		hash := model.Hash(r)

		if _, ok := vehicles[hash]; !ok {
			v := model.NewVehicle(r.Vin, r.Brand, r.Model, r.Year)

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

			vehicles[hash] = &v
		}

		vehicles[hash].AppendRegistrations(r)
	}

	for _, op := range operations {
		logger.Debugf("detectVehicles: operation: %#v", op.String())
		hash := model.Hash(op)

		if _, ok := vehicles[hash]; !ok {
			v := model.NewVehicle(op.Vin, op.Brand, op.Model, op.Year)
			vehicles[hash] = &v
		}

		vehicles[hash].AppendOperations(op)
	}

	return vehicles, nil
}
