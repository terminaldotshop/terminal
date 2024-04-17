package api

type Address struct {
	Name      string  `json:"name"`
	AddrLine1 *string `json:"line1"`
	AddrLine2 *string `json:"line2"`
	City      *string `json:"city"`
	State     *string `json:"state,omitempty"`
	Country   *string `json:"country"`
	Zip       *string `json:"zip"`
}

func NewAddress(name, line1, line2, city, state, country, zip string) Address {
	return Address{
		Name:      name,
		AddrLine1: &line1,
		AddrLine2: &line2,
		City:      &city,
		State:     &state,
		Country:   &country,
		Zip:       &zip,
	}
}
