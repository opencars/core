package customer

import (
	"context"
	"time"

	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/schema"
	"github.com/opencars/seedwork/logger"

	"github.com/opencars/core/pkg/domain"
	"github.com/opencars/core/pkg/domain/model"
	"github.com/opencars/core/pkg/domain/query"
)

type Service struct {
	r domain.RegistrationProvider
	o domain.OperationProvider
	p schema.Producer
}

func NewService(r domain.RegistrationProvider, o domain.OperationProvider, p schema.Producer) (*Service, error) {
	return &Service{
		r: r,
		o: o,
		p: p,
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, q *query.ListByNumber) (*model.Aggregate, error) {
	if err := query.Process(q); err != nil {
		return nil, err
	}

	// Find all registratiton with given number.
	registrations, err := s.r.FindByNumber(ctx, q.Number)
	if err != nil {
		return nil, err
	}

	logger.Debugf("find all operations with given number")

	// Find all operations with given number.
	operations, err := s.o.FindByNumber(ctx, q.Number)
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

		v.AddOpAction(operations...)
		v.AddRegAction(registrations...)
	}

	result := model.NewAggregateWithNumber(q.Number, vehicles)

	if err := s.p.Produce(ctx, q.Event(result)); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) FindByVIN(ctx context.Context, q *query.ListByVIN) (*model.Vehicle, error) {
	if err := query.Process(q); err != nil {
		return nil, err
	}

	// Find all registratiton with given vin.
	registrations, err := s.r.FindByVIN(ctx, q.VIN)
	if err != nil {
		return nil, err
	}

	// Find all operations with given vin.
	operations, err := s.o.FindByVIN(ctx, q.VIN)
	if err != nil {
		return nil, err
	}

	// Detect all unique vehicles from given operations and registrations.
	vehicles, err := s.detectVehicles(ctx, operations, registrations)
	if err != nil {
		return nil, err
	}

	result := model.NewAggregate(vehicles)

	if err := s.p.Produce(ctx, q.Event(result)); err != nil {
		return nil, err
	}

	if len(result.Vehicles) != 0 {
		return &result.Vehicles[0], nil
	}

	return nil, nil
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

		vehicles[hash].AddRegAction(r)
	}

	for _, op := range operations {
		logger.Debugf("detectVehicles: operation: %#v", op.String())
		hash := model.Hash(op)

		if _, ok := vehicles[hash]; !ok {
			v := model.NewVehicle(op.Vin, op.Brand, op.Model, op.Year)
			vehicles[hash] = &v
		}

		vehicles[hash].AddOpAction(op)
	}

	return vehicles, nil
}
