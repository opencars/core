package domain

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"sort"
	"time"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
)

type RegistrationEntity int

const (
	UnknownEntity RegistrationEntity = iota
	Individual
	Legal
)

type Registration struct {
	Code        string
	Number      string
	Vin         string
	Brand       string
	Model       string
	Year        int32
	Capacity    int32
	Color       string
	Fuel        string
	Kind        string
	NumSeating  int32
	OwnWeight   int32
	TotalWeight int32
	Date        *time.Time

	Category string
}

type Department struct {
	Code int32
	Name string
}

type Territory struct {
	Code int32
}

type Owner struct {
	Entity       RegistrationEntity
	Registration *Territory
}

type Operation struct {
	Number      string
	Brand       string
	Model       string
	Year        int32
	Capacity    int32
	Color       string
	Fuel        string
	Kind        string
	Body        string
	Purpose     string
	OwnWeight   int32
	TotalWeight int32
	Date        *time.Time
	Department  *Department
	Owner       *Owner
}

type Vehicle struct {
	VIN          *core.Vin
	FirstRegDate *time.Time
	Brand        string
	Model        string
	Year         int32

	RegistrationExist map[[sha1.Size]byte]struct{}
	OperationExist    map[[sha1.Size]byte]struct{}

	Registrations []*registration.Record
	Operations    []*operation.Record
}

func (v *Vehicle) LastRegistration() *registration.Record {
	if len(v.Registrations) == 0 {
		return nil
	}

	last := v.Registrations[0]

	for i := 1; i < len(v.Registrations); i++ {

		if dateAfterThan(v.Registrations[i].Date, last.Date) {
			last = v.Registrations[i]
		}
	}

	return last
}

func (v *Vehicle) LastOperation() *operation.Record {
	if len(v.Operations) == 0 {
		return nil
	}

	last := v.Operations[0]

	for i := 1; i < len(v.Operations); i++ {

		if dateAfterThan(v.Operations[i].Date, last.Date) {
			last = v.Operations[i]
		}
	}

	return last
}

func (v *Vehicle) LastModificationAt() time.Time {
	o := v.LastOperation()
	r := v.LastRegistration()

	var ot time.Time
	var rt time.Time

	if o.Date != nil {
		ot = time.Date(int(o.Date.Year), time.Month(o.Date.Month), int(o.Date.Day), 0, 0, 0, 0, time.UTC)
	}

	if r.Date != nil {
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

		RegistrationExist: make(map[[sha1.Size]byte]struct{}),
		OperationExist:    make(map[[sha1.Size]byte]struct{}),

		Registrations: make([]*registration.Record, 0),
		Operations:    make([]*operation.Record, 0),
	}
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

		if v.OperationExist == nil {
			v.OperationExist = make(map[[20]byte]struct{})
		}

		_, ok := v.OperationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped", hex)
			continue
		}

		v.OperationExist[sha1] = struct{}{}
		v.Operations = append(v.Operations, candidate)

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

		_, ok := v.RegistrationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped", hex)
			continue
		}

		v.RegistrationExist[sha1] = struct{}{}
		v.Registrations = append(v.Registrations, candidate)

		// Try to assign vin code if it is not already assinged.
		if candidate.Vin != "" && !v.HasVIN() {
			candidate.Vin = v.VIN.Value
		}
	}
}

type Aggregate struct {
	Vehicles []Vehicle
}

func NewAggregate(vehicles map[string]*Vehicle) *Aggregate {
	sorted := make([]Vehicle, 0, len(vehicles))

	// Copy.
	for _, v := range vehicles {
		sorted = append(sorted, *v)
	}

	// Sort vehicles by last modification in operations or registrations.
	sort.Slice(sorted, func(i, j int) bool {
		x := sorted[i].LastModificationAt()
		y := sorted[j].LastModificationAt()

		return y.After(x)
	})

	return &Aggregate{
		Vehicles: sorted,
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

type Hashable interface {
	GetBrand() string
	GetModel() string
	GetYear() int32
	GetCapacity() int32
}

func Hash(x Hashable) string {
	key := fmt.Sprintf("%s-%s-%d-%d", x.GetBrand(), x.GetModel(), x.GetYear(), x.GetCapacity())
	sha1 := sha1.Sum([]byte(key))

	return base64.URLEncoding.EncodeToString(sha1[:])
}

func dateAfterThan(x *common.Date, y *common.Date) bool {
	if x == nil || y == nil {
		return false
	}

	xt := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, time.UTC)
	yt := time.Date(int(y.Year), time.Month(y.Month), int(y.Day), 0, 0, 0, 0, time.UTC)

	return xt.After(yt)
}
