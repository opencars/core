package model

import (
	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/grpc/pkg/operation"
	"github.com/opencars/grpc/pkg/registration"
	"github.com/opencars/grpc/pkg/wanted"
)

type actionSource string

const (
	actionSourceOperation    actionSource = "Operation"
	actionSourceRegistration actionSource = "Registration"
	actionSourceWanted       actionSource = "Wanted"
)

type Action struct {
	source actionSource

	Number      string
	Vin         string
	Brand       string
	Model       string
	Year        int32
	Capacity    int32
	Color       string
	Fuel        string
	Kind        string
	Body        string
	Date        *Date
	OwnWeight   int32
	TotalWeight int32

	// Operation:
	Purpose    *string
	Action     *operation.RecordAction
	Department *operation.Department
	Owner      *operation.Owner

	// Registration:
	Code         *string
	NumSeating   *int32
	FirstRegDate *common.Date
	Category     *string
}

func NewActionFromOperation(dto *operation.Record) *Action {
	return &Action{
		source:      actionSourceOperation,
		Number:      dto.Number,
		Vin:         dto.Vin,
		Brand:       dto.Brand,
		Model:       dto.Model,
		Year:        dto.Year,
		Capacity:    dto.Capacity,
		Color:       dto.Color,
		Fuel:        dto.Fuel,
		Kind:        dto.Kind,
		Body:        dto.Body,
		Date:        NewDateFromProto(dto.Date),
		OwnWeight:   dto.OwnWeight,
		TotalWeight: dto.TotalWeight,

		// Specific:
		Purpose:    &dto.Purpose,
		Action:     dto.Action,
		Department: dto.Department,
		Owner:      dto.Owner,
	}
}

func NewActionFromRegistration(dto *registration.Record) *Action {
	return &Action{
		source:      actionSourceRegistration,
		Number:      dto.Number,
		Vin:         dto.Vin,
		Brand:       dto.Brand,
		Model:       dto.Model,
		Year:        dto.Year,
		Capacity:    dto.Capacity,
		Color:       dto.Color,
		Fuel:        dto.Fuel,
		Kind:        dto.Kind,
		Body:        dto.Body,
		Date:        NewDateFromProto(dto.Date),
		OwnWeight:   dto.OwnWeight,
		TotalWeight: dto.TotalWeight,

		// Specific:
		Code:         &dto.Code,
		NumSeating:   &dto.NumSeating,
		FirstRegDate: dto.FirstRegDate,
		Category:     &dto.Category,
	}
}

func (a *Action) MergeOperation(dto *operation.Record) {
	if a.source == actionSourceOperation {
		return
	}

	// VIN-code from operations has higher priority.
	if dto.Vin != "" {
		a.Vin = dto.Vin
	}

	if a.Body == "" && dto.Body != "" {
		a.Body = dto.Body
	}

	a.Purpose = &dto.Purpose
	a.Action = dto.Action
	a.Department = dto.Department
	a.Owner = dto.Owner
}

func (a *Action) MergeRegistration(dto *registration.Record) {
	if a.source == actionSourceRegistration {
		return
	}

	if a.Vin == "" && dto.Vin != "" {
		a.Vin = dto.Vin
	}

	if a.Body == "" && dto.Body != "" {
		a.Body = dto.Body
	}

	a.Code = &dto.Code
	a.NumSeating = &dto.NumSeating
	a.FirstRegDate = dto.FirstRegDate
	a.Category = &dto.Category
}

func (a *Action) toGRPC() *core.Action {
	dto := core.Action{
		Number:      a.Number,
		Vin:         a.Vin,
		Brand:       a.Brand,
		Model:       a.Model,
		Year:        a.Year,
		Capacity:    a.Capacity,
		Color:       a.Color,
		Fuel:        a.Fuel,
		Kind:        a.Kind,
		Body:        a.Body,
		Date:        a.Date.ToProto(),
		OwnWeight:   a.OwnWeight,
		TotalWeight: a.TotalWeight,
	}

	if a.Purpose != nil {
		dto.Purpose = *a.Purpose
	}

	if a.Action != nil {
		dto.Action = a.Action
	}

	if a.Department != nil {
		dto.Department = a.Department
	}

	if a.Owner != nil {
		dto.Owner = a.Owner
	}

	if a.Code != nil {
		dto.Code = *a.Code
	}

	if a.NumSeating != nil {
		dto.NumSeating = *a.NumSeating
	}

	if a.FirstRegDate != nil {
		dto.FirstRegDate = a.FirstRegDate
	}

	if a.Category != nil {
		dto.Category = *a.Category
	}

	return &dto
}

func NewActionFromWanted(dto *wanted.Vehicle) *Action {
	return &Action{
		source: actionSourceWanted,
		Number: dto.Number,
		Vin:    dto.Vin,
		Brand:  dto.Brand,
		Model:  dto.Model,
		Color:  dto.Color,
		Kind:   dto.Kind,
	}
}
