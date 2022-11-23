package adverts

type requestBody struct {
	VINs    []string `json:"vins"`
	Numbers []string `json:"registration_numbers"`
}

func newRequestBody(vins []string, numbers []string) requestBody {
	return requestBody{
		VINs:    vins,
		Numbers: numbers,
	}
}
