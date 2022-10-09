package domain

import (
	"crypto/sha1"
	"fmt"
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

type Hashable interface {
	GetDate() *common.Date
}

// AppendOperations guarantees uniqness of the operations set.
func (v *Vehicle) AppendOperations(candidates ...*operation.Record) {
	for _, candidate := range candidates {
		candidate.GetDate()
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%s", candidate.Action.Code, date)
		sha1 := sha1.Sum([]byte(s))

		if v.OperationExist == nil {
			v.OperationExist = make(map[[20]byte]struct{})
		}

		_, ok := v.OperationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped")
			continue
		}

		v.OperationExist[sha1] = struct{}{}
		v.Operations = append(v.Operations, candidate)
	}
}

// AppendRegistrations guarantees uniqness of the operations set.
func (v *Vehicle) AppendRegistrations(candidates ...*registration.Record) {
	for _, candidate := range candidates {
		candidate.GetDate()
		date := fmt.Sprintf("%d-%d-%d", candidate.Date.Day, candidate.Date.Month, candidate.Date.Year)
		s := fmt.Sprintf("%d-%d-%s", candidate.Capacity, candidate.Year, date)
		sha1 := sha1.Sum([]byte(s))

		if v.RegistrationExist == nil {
			v.RegistrationExist = make(map[[20]byte]struct{})
		}

		_, ok := v.RegistrationExist[sha1]
		if ok {
			logger.Debugf("candidate %s skipped")
			continue
		}

		v.RegistrationExist[sha1] = struct{}{}
		v.Registrations = append(v.Registrations, candidate)
	}
}

type Aggregate struct {
	Vehicles map[string]*Vehicle
}

func (aggr *Aggregate) VINs() []string {
	vins := make([]string, 0, len(aggr.Vehicles))
	for vin := range aggr.Vehicles {
		vins = append(vins, vin)
	}

	return vins
}
