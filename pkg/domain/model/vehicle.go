package model

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/seedwork/logger"
	"github.com/opencars/translit"
)

type Vehicle struct {
	VIN          *core.Vin
	FirstRegDate *time.Time
	Brand        string
	Model        string
	Year         int32

	registrationExist map[[sha1.Size]byte]struct{}
	operationExist    map[[sha1.Size]byte]struct{}

	registrations  []*registration.Record
	operations     []*operation.Record
	advertisements []Advertisement
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

		registrationExist: make(map[[sha1.Size]byte]struct{}),
		operationExist:    make(map[[sha1.Size]byte]struct{}),

		registrations: make([]*registration.Record, 0),
		operations:    make([]*operation.Record, 0),
	}
}

func (v *Vehicle) LastRegistrationWithNumber(number string) *registration.Record {
	if len(v.registrations) == 0 {
		return nil
	}

	var last *registration.Record
	maxTime := &common.Date{}

	for i := 0; i < len(v.registrations); i++ {
		if translit.ToLatin(number) != translit.ToLatin(v.registrations[i].Number) {
			continue
		}

		if dateAfterThan(v.registrations[i].Date, maxTime) {
			last = v.registrations[i]
			maxTime = last.Date
		}
	}

	return last
}

func (v *Vehicle) LastOperationWithNumber(number string) *operation.Record {
	if len(v.operations) == 0 {
		return nil
	}

	var last *operation.Record
	maxTime := &common.Date{}

	for i := 0; i < len(v.operations); i++ {
		if translit.ToLatin(number) != translit.ToLatin(v.operations[i].Number) {
			continue
		}

		if dateAfterThan(v.operations[i].Date, maxTime) {
			last = v.operations[i]
			maxTime = last.Date
		}
	}

	return last
}

func (v *Vehicle) LastModificationWithNumber(number string) time.Time {
	o := v.LastOperationWithNumber(number)
	r := v.LastRegistrationWithNumber(number)

	var ot time.Time
	var rt time.Time

	if o != nil && o.Date != nil {
		ot = time.Date(int(o.Date.Year), time.Month(o.Date.Month), int(o.Date.Day), 0, 0, 0, 0, time.UTC)
	}

	if r != nil && r.Date != nil {
		rt = time.Date(int(r.Date.Year), time.Month(r.Date.Month), int(r.Date.Day), 0, 0, 0, 0, time.UTC)
	}

	if ot.After(rt) {
		return ot
	}

	return rt
}

func (v *Vehicle) HasVIN() bool {
	return v.VIN != nil && v.VIN.Value != ""
}

func (v *Vehicle) SetFirstRegDate(x time.Time) {
	v.FirstRegDate = &x
}

// AppendOperations guarantees uniqness of the operations set.
func (v *Vehicle) AppendOperations(candidates ...*operation.Record) {
	for _, candidate := range candidates {
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%s", candidate.Action.Code, date)
		sha1 := sha1.Sum([]byte(s))
		hex := base64.URLEncoding.EncodeToString(sha1[:])

		if v.operationExist == nil {
			v.operationExist = make(map[[20]byte]struct{})
		}

		_, ok := v.operationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped", hex)
			continue
		}

		v.operationExist[sha1] = struct{}{}
		v.operations = append(v.operations, candidate)

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

// AppendRegistrations guarantees uniqness of the operations set.
func (v *Vehicle) AppendRegistrations(candidates ...*registration.Record) {
	for _, candidate := range candidates {
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%d-%s", candidate.Capacity, candidate.Year, date)
		sha1 := sha1.Sum([]byte(s))
		hex := base64.URLEncoding.EncodeToString(sha1[:])

		_, ok := v.registrationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped", hex)
			continue
		}

		v.registrationExist[sha1] = struct{}{}
		v.registrations = append(v.registrations, candidate)

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

func (v *Vehicle) AppendAdvertisements(candidates ...Advertisement) {
	v.advertisements = append(v.advertisements, candidates...)
}

func (v *Vehicle) GetOperations() []*operation.Record {
	return v.operations
}

func (v *Vehicle) GetRegistrations() []*registration.Record {
	return v.registrations
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

	dto.Registrations = v.registrations
	dto.Operations = v.operations

	return &dto
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
