package api

import "github.com/stripe/stripe-go/v78"

func NewCardParams(
	number string,
	expMonth string,
	expYear string,
	cvc string,
	line1 string,
	line2 string,
	state string,
	country string,
	city string,
	zip string,
) *stripe.CardParams {
	return &stripe.CardParams{
		Number:         stripe.String(number),
		ExpMonth:       stripe.String(expMonth),
		ExpYear:        stripe.String(expYear),
		CVC:            stripe.String(cvc),
		AddressLine1:   stripe.String(line1),
		AddressLine2:   stripe.String(line2),
		AddressCountry: stripe.String(country),
		AddressCity:    stripe.String(city),
		AddressState:   stripe.String(state),
		AddressZip:     stripe.String(zip),
	}

}
