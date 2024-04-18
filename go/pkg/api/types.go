package api

import (
	"strings"

	"github.com/stripe/stripe-go/v78"
	"github.com/terminalhq/terminal/go/pkg/assert"
)

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
	// If we have an actual card, let's double check a few things
	if name != "" {
		number = strings.ReplaceAll(number, " ", "")

		assert.Assert(len(number) == 16, "Invalid credit card number")
		assert.Assert(len(expMonth) == 2, "Invalid expiration month")
		assert.Assert(len(expYear) == 2, "Invalid expiration year")
	}

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
