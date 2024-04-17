package main

import (
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
}
