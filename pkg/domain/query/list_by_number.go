package query

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/opencars/schema"
	"github.com/opencars/schema/vehicle"
	"github.com/opencars/seedwork"
	"github.com/opencars/translit"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/opencars/core/pkg/domain/model"
)

type ListByNumber struct {
	UserID  string
	TokenID string
	Number  string
}

func (q *ListByNumber) Prepare() {
	q.Number = translit.ToLatin(strings.ToUpper(q.Number))
}

func (q *ListByNumber) Validate() error {
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
			&q.Number,
			validation.Required.Error(seedwork.Required),
			validation.Length(2, 18).Error(seedwork.Invalid),
		),
	)
}

func (q *ListByNumber) Event(result *model.Aggregate) schema.Producable {
	msg := vehicle.VehicleSearched{
		UserId:       q.UserID,
		TokenId:      q.TokenID,
		Number:       q.Number,
		ResultAmount: uint32(len(result.Vehicles)),
		SearchedAt:   timestamppb.New(time.Now().UTC()),
	}

	return schema.New(&source, &msg).Message(
		schema.WithSubject(schema.CoreCustomerActions),
	)
}
