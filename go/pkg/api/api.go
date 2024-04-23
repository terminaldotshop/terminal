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

const apiURL = "https://api.terminal.shop/"
const authURL = "https://auth.terminal.shop/"

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
	product.Art = `
                {
             }   }   {
            {   {  }  }
             }   }{  {
           _{  }{  }  }_
          (  }{  }{  {  )
          |""---------""|
          |             /""\
          |            | _  |
          |             / | |
          |             |/  |
          |             /  /
          |            |  /
          |            "T"
           ""---------""`

	return &product, nil
}

type FingerprintRequest struct {
	Fingerprint string `json:"fingerprint"`
}

type UserCredentials struct {
	// TODO: Get the email as well from this?

	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}

func (u UserCredentials) String() string {
	return fmt.Sprintf("{UserID: '%s', AccessToken: '%s...'}", u.UserID, string(u.AccessToken[:10]))
}

func FetchUserToken(publicKey string) (*UserCredentials, error) {
	fingerprint := FingerprintRequest{Fingerprint: publicKey}

	resp, body, err := post("", "ssh/login", fingerprint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		return nil, errors.New(fmt.Sprintf("Server error: %s", string(body)))
	}

	var creds UserCredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, err
	}

	if creds.AccessToken == "" {
		return nil, errors.New(fmt.Sprintf("Failed to fetch: %s", publicKey))
	}

	return &creds, nil
}

type OrderParams struct {
	Email    string         `json:"email"`
	Shipping Address        `json:"shipping"`
	Products []ProductOrder `json:"products"`
}

type ProductOrder struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

// {"id":"in_1r65m1t9LEE","subtotal":7500,"shipping":1000,"total":8500}
type ShippingInfoResponse struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"name"`
	Cost        float64 `json:"cost"`
	Estimate    string  `json:"estimate"`
}

type OrderResponse struct {
	OrderID  string               `json:"id"`
	Shipping ShippingInfoResponse `json:"shipping"`
	Tax      int                  `json:"tax"`
	Subtotal int                  `json:"subtotal"`
	Total    int                  `json:"total"`
}

func CreateOrder(bearer string, order OrderParams) (*OrderResponse, error) {
	_, body, err := post(bearer, "api/order", order)
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

type SubmitOrderResponse bool

func PurchaseOrder(bearer string, orderID string, cardInfo *StripeCardToken) (*SubmitOrderResponse, error) {
	_, body, err := post(bearer, "api/payment", SubmitOrderRequest{
		OrderID:         orderID,
		StripeCardToken: cardInfo.GetToken(),
	})
	if err != nil {
		return nil, err
	}

	var response SubmitOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(fmt.Sprintf("Error: %s, body: %s", err, string(body)))
	}

	return &response, nil
}

func Init(testing bool) {
	// THIS IS A PUBLIC KEY, SO ITS OK IF PRIME LEAKS THIS. ITS PUBLIC!!
	//	(that's what `pk` stands for - public key)
	//
	//	You can prove it to yourself here: https://sourcegraph.com/search?q=context:global+pk_live_&patternType=keyword&sm=0
	if testing {
		stripe.Key = "pk_test_51OrLKuDgGJQx1Mr6CNDnGNukgQwBonSSToYC8VcmE7qBk2YEad8UuitmY54Pqp0tuZCrk8PNP9cEKVYHvuLcjsnv007CKDgOew"

		return
	}

	stripe.Key = "pk_live_51OrLKuDgGJQx1Mr6tJbUNOAWOcAZ1gGs2rr6T99REuLD6tPPPfSS6iSZnLAI7Kw0EBR63m8SIcqdeb3vhVRLbqZr00zy2bzLav"
}

func post(bearer string, path string, payload any) (*http.Response, []byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	request, err := http.NewRequest("POST", routeAPI(path), bytes.NewReader(buf))
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))
	resp, err := http.DefaultClient.Do(request)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}
