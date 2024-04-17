package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v78"
	// "github.com/stripe/stripe-go/v78/client"
	// "github.com/stripe/stripe-go/v78/paymentintent"
	"github.com/stripe/stripe-go/v78/token"
	// "github.com/stripe/stripe-go/v78/customer"
)

// const apiURL = "https://api.thdxr.dev.terminal.shop/"
// const authURL = "https://auth.thdxr.dev.terminal.shop/"

const apiURL = "https://api.adam.dev.terminal.shop/"
const authURL = "https://auth.adam.dev.terminal.shop/"

func routeAPI(path string) string {
	return fmt.Sprintf("%s%s", apiURL, path)
}

func routeAuth(path string) string {
	return fmt.Sprintf("%s%s", authURL, path)
}

func FetchProducts() (ProductResponse, error) {
	resp, err := http.Get(routeAPI("api/product"))
	if err != nil {
		return nil, err
	}

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productResponse ProductResponse
	if err := json.Unmarshal(body, &productResponse); err != nil {
		return nil, err
	}

	return productResponse, nil
}

func FetchOneProduct() (*Product, error) {
	response, err := FetchProducts()
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("No response somehow....")
	}

	product := response[0]
	return &product, nil
}

type FingerprintRequest struct {
	Fingerprint string `json:"fingerprint"`
}

type UserCredentials struct {
	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}

func (u UserCredentials) String() string {
	return fmt.Sprintf("{UserID: '%s', AccessToken: '%s...'}", u.UserID, string(u.AccessToken[:10]))
}

func FetchUserToken(public_key string) (*UserCredentials, error) {
	fingerprint := FingerprintRequest{Fingerprint: public_key}
	marshalled, _ := json.Marshal(fingerprint)
	resp, err := http.Post(routeAuth("ssh/login"), "application/json", bytes.NewReader(marshalled))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode == 500 {
		return nil, errors.New(fmt.Sprintf("Server error: %s", string(body)))
	}

	var creds UserCredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

type OrderParams struct {
	Email    string          `json:"email"`
	Shipping ShippingDetails `json:"shipping"`
	Products []ProductOrder  `json:"products"`
}

type ShippingDetails struct {
	Name     string  `json:"name"`
	Address1 *string `json:"line1"`
	Address2 *string `json:"line2"`
	City     *string `json:"city"`
	State    *string `json:"state,omitempty"`
	Country  *string `json:"country"`
	Zip      *string `json:"zip"`
}

type ProductOrder struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

// {"id":"in_1P6deNDgGJQx1Mr65m1t9LEE","subtotal":7500,"shipping":1000,"total":8500}
type ShippingInfoResponse struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"name"`
	Cost        float64 `json:"cost"`
	Estimate    string  `json:"estimate"`
}

type OrderResponse struct {
	OrderID  string               `json:"id"`
	Tax      int                  `json:"tax"`
	Shipping ShippingInfoResponse `json:"shipping"`
	Subtotal int                  `json:"subtotal"`
	Total    int                  `json:"total"`
}

func CreateOrder(bearer string, order OrderParams) (*OrderResponse, error) {
	buf, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	// resp, err := http.Post(routeAPI("api/order"), "application/json", bytes.NewReader(buf))
	request, err := http.NewRequest("POST", routeAPI("api/order"), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	resp, err := http.DefaultClient.Do(request)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orderResponse OrderResponse
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return nil, err
	}

	if orderResponse.OrderID == "" {
		return nil, errors.New(fmt.Sprintf("Error: %s, body: %s", "No Order ID", string(body)))
	}

	return &orderResponse, nil
}

type StripeCardToken struct {
	Token *stripe.Token
}

func (s StripeCardToken) GetToken() string {
	// TODO: Not sure which is the right ID to send for Dax
	return s.Token.ID
}

func StripeCreditCard(card *stripe.CardParams) (*StripeCardToken, error) {
	// THIS IS A PUBLISHABLE TOKEN, SO ITS OK IF PRIME LEAKS THIS. ITS PUBLIC
	stripe.Key = "pk_test_51OrLKuDgGJQx1Mr6CNDnGNukgQwBonSSToYC8VcmE7qBk2YEad8UuitmY54Pqp0tuZCrk8PNP9cEKVYHvuLcjsnv007CKDgOew"

	tokenParams := &stripe.TokenParams{Card: card}
	tokenResult, err := token.New(tokenParams)

	if err != nil {
		return nil, err
	}

	info := &StripeCardToken{
		Token: tokenResult,
	}

	return info, nil
}

type SubmitOrderRequest struct {
	OrderID         string `json:"orderID"`
	StripeCardToken string `json:"token"`
}

type SubmitOrderResponse struct {
}

func PurchaseOrder(bearer string, orderID string, cardInfo *StripeCardToken) (*SubmitOrderResponse, error) {
	buf, err := json.Marshal(SubmitOrderRequest{
		OrderID:         orderID,
		StripeCardToken: cardInfo.GetToken(),
	})

	if err != nil {
		return nil, err
	}

	// resp, err := http.Post(routeAPI("api/payment"), "application/json", bytes.NewReader(buf))
	request, err := http.NewRequest("POST", routeAPI("api/payment"), bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	resp, err := http.DefaultClient.Do(request)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("body:", string(body))

	var response SubmitOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(fmt.Sprintf("Error: %s, body: %s", err, string(body)))
	}

	return &response, nil
}
