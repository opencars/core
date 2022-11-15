package model

import (
	"time"

	"github.com/opencars/grpc/pkg/common"
)

func dateAfterThan(x *common.Date, y *common.Date) bool {
	if x == nil || y == nil {
		return false
	}

	xt := time.Date(int(x.Year), time.Month(x.Month), int(x.Day), 0, 0, 0, 0, time.UTC)
	yt := time.Date(int(y.Year), time.Month(y.Month), int(y.Day), 0, 0, 0, 0, time.UTC)

	return xt.After(yt)
}
