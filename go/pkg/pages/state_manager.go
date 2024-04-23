package pages

import (
	"strings"

	"github.com/terminalhq/terminal/go/pkg/api"
)

const (
	goToEmail     = 1
	goToShipping  = 2
	goToCC        = 3
	goToCCAddr    = 4
	goToConfirm   = 5
	goToAnimation = 6
)

var defaultShippingState = api.NewAddress(
	"Teej DV",
	"Streamer Lane",
	"",
	"Miami",
	"FL",
	"US",
	"33131",
)

var defaultCreditCard = api.NewCreditCard(
	"Teej DV",
	"4242 4242 4242 4242",
	"12",
	"34",
	"314",
)

var defaultBillingAddress = api.NewAddress(
	"TJ DeVries",
	"Credit Card Drive",
	"WithSecondLine",
	"Miami",
	"FL",
	"US",
	"33131",
)

var defaultEmail = "teej_dv@twitch.tv"

func stateToNumber(toState string) int {
	switch strings.ToLower(toState) {
	case "email":
		return goToEmail
	case "shipping":
		return goToShipping
	case "cc":
		return goToCC
	case "cc-addr":
		return goToCCAddr
	case "confirm":
		return goToConfirm
	case "animation":
		return goToAnimation
	}
	return 0
}
