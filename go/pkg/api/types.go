package api

import "github.com/stripe/stripe-go/v78"

type Address struct {
	Name      string `json:"name"`
	AddrLine1 string `json:"line1"`
	AddrLine2 string `json:"line2"`
	City      string `json:"city"`
	State     string `json:"state,omitempty"`
	Country   string `json:"country"`
	Zip       string `json:"zip"`
}

func EmptyAddress() Address {
	return Address{}
}

func NewAddress(name, line1, line2, city, state, country, zip string) Address {
	return Address{
		Name:      name,
		AddrLine1: line1,
		AddrLine2: line2,
		City:      city,
		State:     state,
		Country:   country,
		Zip:       zip,
	}
}

type CreditCard struct {
	Name     string
	Number   string
	ExpMonth string
	ExpYear  string
	CVC      string
}

func NewCreditCard(name, number, expMonth, expYear, cvc string) CreditCard {
	return CreditCard{
		Name:     name,
		Number:   number,
		ExpMonth: expMonth,
		ExpYear:  expYear,
		CVC:      cvc,
	}
}

func NewDefaultCreditCard() CreditCard {
	return NewCreditCard("", "", "", "", "")
}

func NewCardParams(
	cc CreditCard,
	address Address,
) *stripe.CardParams {
	return &stripe.CardParams{
		Name:           stripe.String(cc.Name),
		Number:         stripe.String(cc.Number),
		ExpMonth:       stripe.String(cc.ExpMonth),
		ExpYear:        stripe.String(cc.ExpYear),
		CVC:            stripe.String(cc.CVC),
		AddressLine1:   stripe.String(address.AddrLine1),
		AddressLine2:   stripe.String(address.AddrLine2),
		AddressCountry: stripe.String(address.Country),
		AddressCity:    stripe.String(address.City),
		AddressState:   stripe.String(address.State),
		AddressZip:     stripe.String(address.Zip),
	}

}
