package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVehicle_AddAction(t *testing.T) {
	v := Vehicle{
		actions: []*Action{
			{
				Number:   "AI8474MK",
				Vin:      "3VWC57BU9KM078586",
				Brand:    "VOLKSWAGEN",
				Model:    "JETTA",
				Year:     2018,
				Capacity: 2000,
				Color:    "Сірий",
				Fuel:     "Бензин",
				Kind:     "Легковий",
				Body:     "Седан",
				Date: &Date{
					Day:   17,
					Month: 03,
					Year:  2021,
				},
				OwnWeight:   1432,
				TotalWeight: 1432,
			},
			{
				Number:   "AI8474MK",
				Vin:      "3VWC57BU9KM078586",
				Brand:    "VOLKSWAGEN",
				Model:    "JETTA",
				Year:     2018,
				Capacity: 2000,
				Color:    "Сірий",
				Fuel:     "Бензин",
				Kind:     "Легковий",
				Body:     "Седан",
				Date: &Date{
					Day:   15,
					Month: 03,
					Year:  2020,
				},
				OwnWeight:   1432,
				TotalWeight: 1432,
			},
		},
	}

	newAction := Action{
		Number:   "AI8474MK",
		Vin:      "3VWC57BU9KM078586",
		Brand:    "VOLKSWAGEN",
		Model:    "JETTA",
		Year:     2018,
		Capacity: 2000,
		Color:    "Сірий",
		Fuel:     "Бензин",
		Kind:     "Легковий",
		Body:     "Седан",
		Date: &Date{
			Day:   01,
			Month: 01,
			Year:  2022,
		},
		OwnWeight:   1432,
		TotalWeight: 1432,
	}

	v.AddAction(&newAction)

	expected := []*Action{
		{
			Number:   "AI8474MK",
			Vin:      "3VWC57BU9KM078586",
			Brand:    "VOLKSWAGEN",
			Model:    "JETTA",
			Year:     2018,
			Capacity: 2000,
			Color:    "Сірий",
			Fuel:     "Бензин",
			Kind:     "Легковий",
			Body:     "Седан",
			Date: &Date{
				Day:   01,
				Month: 01,
				Year:  2022,
			},
			OwnWeight:   1432,
			TotalWeight: 1432,
		},
		{
			Number:   "AI8474MK",
			Vin:      "3VWC57BU9KM078586",
			Brand:    "VOLKSWAGEN",
			Model:    "JETTA",
			Year:     2018,
			Capacity: 2000,
			Color:    "Сірий",
			Fuel:     "Бензин",
			Kind:     "Легковий",
			Body:     "Седан",
			Date: &Date{
				Day:   17,
				Month: 03,
				Year:  2021,
			},
			OwnWeight:   1432,
			TotalWeight: 1432,
		},
		{
			Number:   "AI8474MK",
			Vin:      "3VWC57BU9KM078586",
			Brand:    "VOLKSWAGEN",
			Model:    "JETTA",
			Year:     2018,
			Capacity: 2000,
			Color:    "Сірий",
			Fuel:     "Бензин",
			Kind:     "Легковий",
			Body:     "Седан",
			Date: &Date{
				Day:   15,
				Month: 03,
				Year:  2020,
			},
			OwnWeight:   1432,
			TotalWeight: 1432,
		},
	}

	assert.EqualValues(t, expected, v.actions)
}
