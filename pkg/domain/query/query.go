package query

import (
	"github.com/opencars/schema"
	"github.com/opencars/seedwork"
)

var (
	source = schema.Source{
		Service: "core",
		Version: "1.0",
	}
)

func Process(q seedwork.Query) error {
	q.Prepare()

	return seedwork.Validate(q, "request")
}
