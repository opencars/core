package model

import (
	"strings"

	"github.com/opencars/grpc/pkg/core"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gearbox string

const (
	GearboxAutomatic Gearbox = "Automatic"
	GearboxManual    Gearbox = "Manual"
	GearboxManumatic Gearbox = "Manumatic"
	GearboxVariator  Gearbox = "Variator"
	GearboxAMT       Gearbox = "Automated manual transmission (Robot)"
	GearboxOther     Gearbox = "Other"
)

func (g Gearbox) grpc() core.Advertisement_GearboxType {
	switch g {
	case GearboxAutomatic:
		return core.Advertisement_GEARBOX_AUTOMATIC
	case GearboxManual:
		return core.Advertisement_GEARBOX_MANUAL
	case GearboxManumatic:
		return core.Advertisement_GEARBOX_MANUMATIC
	case GearboxVariator:
		return core.Advertisement_GEARBOX_VARIATOR
	case GearboxAMT:
		return core.Advertisement_GEARBOX_AMT
	case GearboxOther:
		return core.Advertisement_GEARBOX_OTHER
	default:
		return core.Advertisement_GEARBOX_UKNOWN
	}
}

type WheelDriveType string

const (
	FrontWheel WheelDriveType = "Front-wheel"
	RearWheel  WheelDriveType = "Real-wheel"
	AllWheel   WheelDriveType = "All-wheel"
)

func (w WheelDriveType) grpc() core.Advertisement_WheelDriveType {
	switch w {
	case FrontWheel:
		return core.Advertisement_WHEELDRIVE_FRONT
	case RearWheel:
		return core.Advertisement_WHEELDRIVE_REAR
	case AllWheel:
		return core.Advertisement_WHEELDRIVE_ALL
	default:
		return core.Advertisement_WHEELDRIVE_UNKNOWN
	}
}

type Category string

const (
	CategoryCar     Category = "Car"
	CategoryTruck   Category = "Truck"
	CategoryMoto    Category = "Motorcycle"
	CategoryBus     Category = "Bus"
	CategoryWater   Category = "Watercraft"
	CategoryAir     Category = "Aircraft"
	CategoryCamper  Category = "Camper"
	CategoryTrailer Category = "Trailer"
	CategorySpecial Category = "Special vehicles"
)

func (c Category) grpc() core.Advertisement_Category {
	switch c {
	case CategoryCar:
		return core.Advertisement_CATEGORY_CAR
	case CategoryTruck:
		return core.Advertisement_CATEGORY_TRUCK
	case CategoryMoto:
		return core.Advertisement_CATEGORY_MOTO
	case CategoryBus:
		return core.Advertisement_CATEGORY_BUS
	case CategoryWater:
		return core.Advertisement_CATEGORY_WATER
	case CategoryAir:
		return core.Advertisement_CATEGORY_AIR
	case CategoryCamper:
		return core.Advertisement_CATEGORY_CAMPER
	case CategoryTrailer:
		return core.Advertisement_CATEGORY_TRAILER
	case CategorySpecial:
		return core.Advertisement_CATEGORY_SPECIAL
	default:
		return core.Advertisement_CATEGORY_UNKNOWN
	}
}

type FuelType string

const (
	Gasoline    FuelType = "Gasoline / petrol"
	Diesel      FuelType = "Diesel"
	Gas         FuelType = "Gas"
	GasGosoline FuelType = "Gas / gasoline"
	Hybrid      FuelType = "Hybrid"
	Electric    FuelType = "Electric battery"
	Propane     FuelType = "Propane-butane / Liquefied petroleum gas (LPG)"
	Methane     FuelType = "Compressed neutral gas (CNG)"
	Other       FuelType = "Other (biodiesel, ethanol)"
)

func (f FuelType) grpc() core.Advertisement_FuelType {
	switch f {
	case Gasoline:
		return core.Advertisement_FUEL_GASOLINE
	case Diesel:
		return core.Advertisement_FUEL_DIESEL
	case Gas:
		return core.Advertisement_FUEL_GAS
	case GasGosoline:
		return core.Advertisement_FUEL_GAS_GASOLINE
	case Hybrid:
		return core.Advertisement_FUEL_HYBRID
	case Electric:
		return core.Advertisement_FUEL_ELECTRIC
	case Propane:
		return core.Advertisement_FUEL_PROPANE
	case Methane:
		return core.Advertisement_FUEL_METHANE
	case Other:
		return core.Advertisement_FUEL_OTHER
	default:
		return core.Advertisement_FUEL_UNKNOWN
	}
}

type BodyType string

const (
	BodyTypeSedan       BodyType = "Sedan / saloon"
	BodyTypeCrossover   BodyType = "Crossover"
	BodyTypeMinivan     BodyType = "Minivan"
	BodyTypeHatchback   BodyType = "Hatchback"
	BodyTypeWagon       BodyType = "Station wagon / estate / universal"
	BodyTypeCoupe       BodyType = "Coupe"
	BodyTypeConvertible BodyType = "Convertible / cabriolet"
	BodyTypePickup      BodyType = "Pickup"
	BodyTypeLimousine   BodyType = "Limousine"
	BodyTypeLightTruck  BodyType = "Light truck"
	BodyTypeOther       BodyType = "Other"
)

func (b BodyType) grpc() core.Advertisement_BodyType {
	switch b {
	case BodyTypeSedan:
		return core.Advertisement_BODY_SEDAN
	case BodyTypeCrossover:
		return core.Advertisement_BODY_CROSSOVER
	case BodyTypeMinivan:
		return core.Advertisement_BODY_MINIVAN
	case BodyTypeHatchback:
		return core.Advertisement_BODY_HATCHBACK
	case BodyTypeWagon:
		return core.Advertisement_BODY_WAGON
	case BodyTypeCoupe:
		return core.Advertisement_BODY_COUPE
	case BodyTypeConvertible:
		return core.Advertisement_BODY_CONVERTIBLE
	case BodyTypePickup:
		return core.Advertisement_BODY_PICKUP
	case BodyTypeLimousine:
		return core.Advertisement_BODY_LIMOUSINE
	case BodyTypeLightTruck:
		return core.Advertisement_BODY_LIGHTTRUCK
	case BodyTypeOther:
		return core.Advertisement_BODY_OTHER
	default:
		return core.Advertisement_BODY_UNKNOWN
	}
}

type Advertisement struct {
	ID                     int32          `json:"id"`
	ResourceID             string         `json:"resource_id"`
	Brand                  string         `json:"brand_id"`
	Model                  string         `json:"model_id"`
	Category               Category       `json:"category"`
	Title                  string         `json:"title"`
	IsVerified             bool           `json:"is_verified"`
	URL                    string         `json:"url"`
	Price                  int            `json:"price"`
	ImageUrls              []string       `json:"image_urls"`
	LastSeenAt             int            `json:"last_seen_at"`
	PublishedAt            int            `json:"published_at"`
	ScrapedAt              int            `json:"scraped_at"`
	UpdatedAt              int            `json:"updated_at"`
	Year                   int32          `json:"year"`
	Gearbox                Gearbox        `json:"gearbox_type"`
	WheelDrive             WheelDriveType `json:"wheel_drive_type"`
	EngineCapacity         float64        `json:"engine_capacity"`
	Fuel                   FuelType       `json:"fuel_type"`
	Mileage                int            `json:"mileage"`
	Body                   BodyType       `json:"bodystyle"`
	IsCustomsCleared       bool           `json:"is_customs_cleared"`
	VinPage                string         `json:"vin_page"`
	VinOpencars            string         `json:"vin_opencars"`
	RegistrationNumberPage string         `json:"registration_number_page"`
	IDOnResource           int            `json:"id_on_resource"`
}

func (a *Advertisement) toGRPC() *core.Advertisement {
	var vin string
	if a.VinPage != "" {
		vin = a.VinPage
	} else if a.VinOpencars != "" {
		vin = a.VinOpencars
	}

	return &core.Advertisement{
		Id:         uint32(a.ID),
		Resource:   a.ResourceID,
		Brand:      a.Brand,
		Model:      a.Model,
		Category:   a.Category.grpc(),
		Title:      a.Title,
		IsVerified: a.IsVerified,
		Url:        a.URL,
		Price:      uint32(a.Price),
		ImageUrls:  a.ImageUrls,
		LastSeenAt: &timestamppb.Timestamp{
			Seconds: int64(a.LastSeenAt),
		},
		PublishedAt: &timestamppb.Timestamp{
			Seconds: int64(a.PublishedAt),
		},
		ScrapedAt: &timestamppb.Timestamp{
			Seconds: int64(a.ScrapedAt),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: int64(a.UpdatedAt),
		},
		Year:             uint32(a.Year),
		Gearbox:          a.Gearbox.grpc(),
		WheelDrive:       a.WheelDrive.grpc(),
		EngineCapacity:   a.EngineCapacity,
		Fuel:             a.Fuel.grpc(),
		Mileage:          uint32(a.Mileage),
		Body:             a.Body.grpc(),
		IsCustomsCleared: a.IsCustomsCleared,
		Vin:              vin,
		Number:           a.RegistrationNumberPage,
	}
}

func (a *Advertisement) GetBrand() string {
	return strings.ToUpper(a.Brand)
}
func (a *Advertisement) GetModel() string {
	return strings.ToUpper(a.Model)
}
func (a *Advertisement) GetYear() int32 {
	return a.Year
}

func (a *Advertisement) GetCapacity() int32 {
	return int32(a.EngineCapacity * float64(1000))
}
