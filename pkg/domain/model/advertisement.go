package model

type Advertisement struct {
	ID                     int      `json:"id"`
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
