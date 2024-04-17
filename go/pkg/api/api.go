package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func FetchWidgets() (*WidgetResponse, error) {
	resp, err := (http.Get("https://api.terminal.shop/api/product"))
	if err != nil {
		return nil, err
	}

	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var widgetResponse WidgetResponse
	if err := json.Unmarshal(body, &widgetResponse); err != nil {
		return nil, err
	}

	return &widgetResponse, nil
}

func FetchOneProduct() (*Widget, error) {
	response, err := FetchWidgets()
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("No response somehow....")
	}

	widget := response.Results[0]
	return &widget, nil
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
	resp, err := http.Post("https://auth.terminal.shop/ssh/login", "application/json", bytes.NewReader(marshalled))
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

type Order struct {
	Shipping ShippingDetails `json:"shipping"`
	Products []ProductOrder  `json:"products"`
}

type ShippingDetails struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Name     string `json:"name"`
}

type ProductOrder struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type OrderResposne struct{}

func PlaceOrder(order Order) (*OrderResposne, error) {
	buf, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.terminal.shop/api/order", "application/json", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var orderResponse OrderResposne
	if err := json.Unmarshal(body, &orderResponse); err != nil {
		return nil, err
	}

	return &orderResponse, nil

}
