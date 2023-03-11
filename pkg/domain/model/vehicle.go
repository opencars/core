package model

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"time"

	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/core/customer"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/grpc/pkg/wanted"
	"github.com/opencars/seedwork/logger"
)

type Vehicle struct {
	VIN          *core.Vin
	FirstRegDate *time.Time
	Brand        string
	Model        string
	Year         int32

	actionExist map[[sha1.Size]byte]*Action

	wanted         []*wanted.Vehicle
	advertisements []Advertisement
	actions        []*Action
}

func NewVehicle(vin, brand, model string, year int32) Vehicle {
	var vinCode *core.Vin
	if vin != "" {
		vinCode = &core.Vin{Value: vin}
	}

	return Vehicle{
		VIN:   vinCode,
		Brand: brand,
		Model: model,
		Year:  year,

		actionExist: make(map[[sha1.Size]byte]*Action),

		wanted: make([]*wanted.Vehicle, 0),

		advertisements: make([]Advertisement, 0),
		actions:        make([]*Action, 0),
	}
}

func (v *Vehicle) LastActionWithNumber(number string) *Action {
	if len(v.actions) == 0 {
		return nil
	}

	for i, action := range v.actions {
		if action.Number == number {
			return v.actions[i]
		}
	}

	return nil
}

func (v *Vehicle) LastModificationWithNumber(number string) time.Time {
	o := v.LastActionWithNumber(number)
	if o == nil {
		return time.Time{}
	}

	return time.Date(int(o.Date.Year), time.Month(o.Date.Month), int(o.Date.Day), 0, 0, 0, 0, time.UTC)
}

func (v *Vehicle) HasVIN() bool {
	return v.VIN != nil && v.VIN.Value != ""
}

func (v *Vehicle) SetFirstRegDate(x time.Time) {
	v.FirstRegDate = &x
}

func (v *Vehicle) AddRegAction(candidates ...*registration.Record) {
	for _, candidate := range candidates {
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%d-%s", candidate.Capacity, candidate.Year, date)
		sha1 := sha1.Sum([]byte(s))

		if action, ok := v.actionExist[sha1]; ok {
			action.MergeRegistration(candidate)
		} else {
			// Add new.
			newAction := NewActionFromRegistration(candidate)
			v.AddAction(newAction)
			v.actionExist[sha1] = newAction
			continue
		}

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

func (v *Vehicle) AddOpAction(candidates ...*operation.Record) {
	for _, candidate := range candidates {
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%d-%s", candidate.Capacity, candidate.Year, date)
		sha1 := sha1.Sum([]byte(s))

		if action, ok := v.actionExist[sha1]; ok {
			action.MergeOperation(candidate)
		} else {
			// Add new.
			newAction := NewActionFromOperation(candidate)
			v.AddAction(newAction)
			v.actionExist[sha1] = newAction
			continue
		}

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

func (v *Vehicle) AddWantedAction(candidates ...*wanted.Vehicle) {
	for _, candidate := range candidates {
		newAction := NewActionFromWanted(candidate)
		v.AddAction(newAction)

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

func (v *Vehicle) AppendAdvertisements(candidates ...Advertisement) {
	for _, add := range candidates {
		i := sort.Search(len(v.advertisements), func(i int) bool {
			return v.advertisements[i].PublishedAt < add.PublishedAt
		})

		v.advertisements = insertAt(v.advertisements, i, add)
	}
}

func (v *Vehicle) AppendWanted(candidates ...*wanted.Vehicle) {
	for _, wanted := range candidates {
		i := sort.Search(len(v.wanted), func(i int) bool {
			return dateAfterThan(wanted.InsertDate, v.wanted[i].InsertDate)
		})

		v.wanted = insertAt(v.wanted, i, wanted)
	}
}

func (v *Vehicle) AddAction(action *Action) {
	i := sort.Search(len(v.actions), func(i int) bool {
		return !v.actions[i].Date.After(action.Date)
	})

	v.actions = insertAt(v.actions, i, action)
}

func (v *Vehicle) ToGRPC() *core.Vehicle {
	dto := core.Vehicle{
		Vin:   v.VIN,
		Brand: v.Brand,
		Model: v.Model,
		Year:  v.Year,
	}

	if v.FirstRegDate != nil {
		dto.FirstRegDate = &common.Date{
			Year:  int32(v.FirstRegDate.Year()),
			Month: int32(v.FirstRegDate.Month()),
			Day:   int32(v.FirstRegDate.Day()),
		}
	}

	dto.Wanted = v.wanted

	logger.Debugf("to grpc: wanted: %s", v.wanted)

	dto.Advertisements = make([]*core.Advertisement, 0)
	dto.Actions = make([]*core.Action, 0)

	for _, add := range v.advertisements {
		dto.Advertisements = append(dto.Advertisements, add.toGRPC())
	}

	for _, action := range v.actions {
		dto.Actions = append(dto.Actions, action.toGRPC())
	}

	return &dto
}

func (v *Vehicle) ToCustomerGRPC() *customer.Vehicle {
	if v == nil {
		return nil
	}

	dto := customer.Vehicle{
		Brand: v.Brand,
		Model: v.Model,
		Year:  v.Year,
	}

	if v.VIN != nil && v.VIN.Value != "" {
		dto.Vin = v.VIN.Value
	}

	dto.Actions = make([]*core.Action, 0)
	for _, action := range v.actions {
		dto.Actions = append(dto.Actions, action.toGRPC())
		logger.Infof("add: %+v", action.toGRPC())
	}

	return &dto
}

func (v *Vehicle) AssignVIN(vin string) {
	v.VIN = &core.Vin{
		Value: vin,
	}
}

func GetVINs(vehicles map[string]*Vehicle) []string {
	vins := make([]string, 0, len(vehicles))
	for _, v := range vehicles {
		if v.VIN.GetValue() == "" {
			continue
		}

		vins = append(vins, v.VIN.GetValue())
	}

	return vins
}
