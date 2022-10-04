package domain

import (
	"time"

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
	VIN          *core.Vin  `json:"vin"`
	FirstRegDate *time.Time `json:"first_reg_date"`
	Brand        string     `json:"brand"`
	Model        string     `json:"model"`
	Year         int32      `json:"year"`

	Registrations []*registration.Record `json:"registrations"`
	Operations    []*operation.Record    `json:"operations"`
}

type Aggregate struct {
	Vehicles map[string]*Vehicle `json:"vehicles"`
}
