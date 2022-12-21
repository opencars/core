package model

import (
	"sort"

	"github.com/opencars/grpc/pkg/core"
)

type Aggregate struct {
	Vehicles []Vehicle
}

func NewAggregateWithNumber(number string, vehicles map[string]*Vehicle) *Aggregate {
	sorted := make([]Vehicle, 0, len(vehicles))

	// Copy.
	for _, v := range vehicles {
		sorted = append(sorted, *v)
	}

	// Sort vehicles by last modification in operations or registrations.
	sort.Slice(sorted, func(i, j int) bool {
		x := sorted[i].LastModificationWithNumber(number)
		y := sorted[j].LastModificationWithNumber(number)

		return x.After(y)
	})

	return &Aggregate{
		Vehicles: sorted,
	}
}

func NewAggregate(vehicles map[string]*Vehicle) *Aggregate {
	items := make([]Vehicle, 0, len(vehicles))

	// Copy.
	for _, v := range vehicles {
		items = append(items, *v)
	}

	return &Aggregate{
		Vehicles: items,
	}
}

func (a *Aggregate) ToGRPC() *core.Result {
	vehicles := make([]*core.Vehicle, 0, len(a.Vehicles))
	for _, v := range a.Vehicles {
		vehicles = append(vehicles, v.ToGRPC())
	}

	return &core.Result{
		Vehicles: vehicles,
	}
}

func (a *Aggregate) ToExternalGRPC() *core.ResultForCustomer {
	vehicles := make([]*core.ResultForCustomer_Vehicle, 0, len(a.Vehicles))
	for _, v := range a.Vehicles {
		vehicles = append(vehicles, v.ToExternalGRPC())
	}

	return &core.ResultForCustomer{
		Vehicles: vehicles,
	}
}
