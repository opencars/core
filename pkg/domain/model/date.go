package model

import (
	"time"

	"github.com/opencars/grpc/pkg/common"
)

type Date struct {
	Day   int32
	Month int32
	Year  int32
}

func NewDateFromProto(date *common.Date) *Date {
	return &Date{
		Day:   date.Day,
		Month: date.Month,
		Year:  date.Year,
	}
}

func (d *Date) After(x *Date) bool {
	if d == nil || x == nil {
		return false
	}

	dt := time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
	xt := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, time.UTC)

	// TODO: Make own comparison without usage of time package.
	return dt.After(xt)
}

func (d *Date) ToProto() *common.Date {
	return &common.Date{
		Day:   d.Day,
		Month: d.Month,
		Year:  d.Year,
	}
}

// TODO: Remove!
func dateAfterThan(x *common.Date, y *common.Date) bool {
	if x == nil || y == nil {
		return false
	}

	xt := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, time.UTC)
	yt := time.Date(int(y.Year), time.Month(y.Month), int(y.Day), 0, 0, 0, 0, time.UTC)

	return xt.After(yt)
}
