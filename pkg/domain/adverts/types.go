package adverts

type requestBody struct {
	VINs []string `json:"vins"`
}

func newRequestBody(vins ...string) requestBody {
	return requestBody{
		VINs: vins,
	}
}
