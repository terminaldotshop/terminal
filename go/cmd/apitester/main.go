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

	token, err := api.FetchUserToken("testing-2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("token: %+v\n", token)

	card := api.NewCardParams(
		"4242424242424242",
		"12",
		"2034",
		"314",
		"Test Line 1",
		"Test Line 2",
		"FL",
		"US",
		"Miami",
		"33131",
	)

	brazilCard := api.NewCardParams(
		"4242424242424242",
		"12",
		"2034",
		"314",
		"Av. Paulista",
		"",
		"SP",
		"BR",
		"São Paulo",
		"01310-000",
	)
	if false {
		card = brazilCard
	}

	berlinCard := api.NewCardParams(
		"4242424242424242",
		"12",
		"2034",
		"314",
		"Am Kupfergraben 5",
		"",
		"",
		"DE",
		"Berlin",
		"10117",
	)
	if false {
		card = berlinCard
		card.AddressState = nil
	}

	fmt.Println("Requesting order...")
	order, err := api.CreateOrder(token.AccessToken, api.OrderParams{
		Email: "KEKW_ADAM_IS_VEGAN@gmail.com",
		Shipping: api.ShippingDetails{
			Address1: card.AddressLine1,
			Address2: card.AddressLine2,
			City:     card.AddressCity,
			State:    card.AddressState,
			Country:  card.AddressCountry,
			Zip:      card.AddressZip,
			Name:     "teej-dv",
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
