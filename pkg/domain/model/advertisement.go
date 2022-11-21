package model

import "github.com/opencars/grpc/pkg/core"

type Advertisement struct {
	ID                     int32    `json:"id"`
	ResourceID             string   `json:"resource_id"`
	BrandID                string   `json:"brand_id"`
	ModelID                string   `json:"model_id"`
	Category               string   `json:"category"`
	Title                  string   `json:"title"`
	IsVerified             bool     `json:"is_verified"`
	URL                    string   `json:"url"`
	Price                  int      `json:"price"`
	ImageUrls              []string `json:"image_urls"`
	LastSeenAt             int      `json:"last_seen_at"`
	PublishedAt            int      `json:"published_at"`
	ScrapedAt              int      `json:"scraped_at"`
	UpdatedAt              int      `json:"updated_at"`
	Year                   int      `json:"year"`
	GearboxType            string   `json:"gearbox_type"`
	WheelDriveType         string   `json:"wheel_drive_type"`
	EngineCapacity         float64  `json:"engine_capacity"`
	FuelType               string   `json:"fuel_type"`
	Mileage                int      `json:"mileage"`
	Bodystyle              string   `json:"bodystyle"`
	IsCustomsCleared       bool     `json:"is_customs_cleared"`
	VinPage                string   `json:"vin_page"`
	VinOpencars            string   `json:"vin_opencars"`
	RegistrationNumberPage string   `json:"registration_number_page"`
	IDOnResource           int      `json:"id_on_resource"`
}

func (a *Advertisement) toGRPC() *core.Advertisement {
	category := core.Advertisement_CATEGORY_UNKNOWN
	switch a.Category {
	case "Car":
		category = core.Advertisement_CATEGORY_CAR
	case "Truck":
		category = core.Advertisement_CATEGORY_TRUCK
	case "Motorcycle":
		category = core.Advertisement_CATEGORY_MOTO
	case "Bus":
		category = core.Advertisement_CATEGORY_BUS
	case "Watercraft":
		category = core.Advertisement_CATEGORY_WATER
	case "Aircraft":
		category = core.Advertisement_CATEGORY_AIR
	case "Camper":
		category = core.Advertisement_CATEGORY_CAMPER
	case "Trailer":
		category = core.Advertisement_CATEGORY_TRAILER
	case "Special vehicles":
		category = core.Advertisement_CATEGORY_SPECIAL
	}

	gearbox := core.Advertisement_GEARBOX_UKNOWN
	switch a.GearboxType {
	case "Automatic":
		gearbox = core.Advertisement_GEARBOX_AUTOMATIC
	case "Manual":
		gearbox = core.Advertisement_GEARBOX_MANUAL
	case "Manumatic":
		gearbox = core.Advertisement_GEARBOX_MANUMATIC
	case "Variator":
		gearbox = core.Advertisement_GEARBOX_VARIATOR
	case "Automated manual transmission (Robot)":
		gearbox = core.Advertisement_GEARBOX_AMT
	case "Other":
		gearbox = core.Advertisement_GEARBOX_OTHER
	}

	return &core.Advertisement{
		Id:          uint32(a.ID),
		Resource:    a.ResourceID,
		Brand:       a.BrandID,
		Model:       a.ModelID,
		Category:    category,
		Title:       a.Title,
		IsVerified:  a.IsVerified,
		Url:         a.URL,
		Price:       uint32(a.Price),
		ImageUrls:   a.ImageUrls,
		LastSeenAt:  uint32(a.LastSeenAt),
		PublishedAt: uint32(a.PublishedAt),
		ScrapedAt:   uint32(a.ScrapedAt),
		UpdatedAt:   uint32(a.UpdatedAt),
		Year:        uint32(a.Year),
		Gearbox:     gearbox,
		// WheelDrive:             a.WheelDrive,
		EngineCapacity: a.EngineCapacity,
		// Fuel:                   a.Fuel,
		Mileage: uint32(a.Mileage),
		// Body:                   a.Body,
		IsCustomsCleared: a.IsCustomsCleared,
		// Vin:                    a.VI,
		RegistrationNumberPage: a.RegistrationNumberPage,
		IdOnResource:           uint32(a.IDOnResource),
	}
}
