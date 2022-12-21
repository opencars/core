package query

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/opencars/schema"
	"github.com/opencars/schema/vehicle"
	"github.com/opencars/seedwork"
	"github.com/opencars/translit"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/opencars/core/pkg/domain/model"
)

type ListByVIN struct {
	UserID  string
	TokenID string
	VIN     string
}

func (q *ListByVIN) Prepare() {
	q.VIN = translit.ToLatin(strings.ToUpper(q.VIN))
}

func (q *ListByVIN) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(
			&q.UserID,
			validation.Required.Error(seedwork.Required),
		),
		validation.Field(
			&q.TokenID,
			validation.Required.Error(seedwork.Required),
		),
		validation.Field(
			&q.VIN,
			validation.Required.Error(seedwork.Required),
			validation.Length(6, 18).Error(seedwork.Invalid),
		),
	)
}

func (q *ListByVIN) Event(result *model.Aggregate) schema.Producable {
	msg := vehicle.RegistrationSearched{
		UserId:       q.UserID,
		TokenId:      q.TokenID,
		Vin:          q.VIN,
		ResultAmount: uint32(len(result.Vehicles)),
		SearchedAt:   timestamppb.New(time.Now().UTC()),
	}

	return schema.New(&source, &msg).Message(
		schema.WithSubject(schema.CoreCustomerActions),
	)
}
