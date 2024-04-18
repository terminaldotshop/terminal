package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/terminalhq/terminal/go/pkg/api"
)

func main() {
	product, err := api.FetchOneProduct()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("product: %s\n", product.Name)

	publicKey := "testing-2"
	token, err := api.FetchUserToken(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("token: %+v\n", token)

	card := api.NewCardParams(
		api.NewCreditCard(
			"Miami User",
			"4242424242424242",
			"12",
			"34",
			"314",
		),
		api.NewAddress(
			"Miami User",
			"Test Line 1",
			"Test Line 2",
			"FL",
			"US",
			"Miami",
			"33131",
		),
	)

	brazilCard := api.NewCardParams(
		api.NewCreditCard(
			"Brazil User",
			"4242424242424242",
			"12",
			"34",
			"314",
		),
		api.NewAddress(
			"Brazil User",
			"Av. Paulista",
			"",
			"SP",
			"BR",
			"São Paulo",
			"01310-000",
		),
	)
	if false {
		card = brazilCard
	}

	fmt.Println("Requesting order...")
	order, err := api.CreateOrder(token.AccessToken, api.OrderParams{
		Email: "KEKW_ADAM_IS_VEGAN@gmail.com",
		Shipping: api.Address{
			AddrLine1: *card.AddressLine1,
			AddrLine2: *card.AddressLine2,
			City:      *card.AddressCity,
			State:     *card.AddressState,
			Country:   *card.AddressCountry,
			Zip:       *card.AddressZip,
			Name:      "teej-dv",
		},
		Products: []api.ProductOrder{
			{
				ID:       product.ID,
				Quantity: 3,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	pretty, _ := json.MarshalIndent(order, "", "  ")
	fmt.Println("order:", string(pretty))

	fmt.Println("Requesting card info...")
	info, err := api.StripeCreditCard(card)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Purchasing Order...")
	orderResp, err := api.PurchaseOrder(token.AccessToken, order.OrderID, info)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("orderResp: %+v\n", orderResp)
}
