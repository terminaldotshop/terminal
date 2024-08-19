package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stripe/stripe-go/v78"
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal/go/pkg/resource"

	"github.com/stripe/stripe-go/v78/token"
)

func Init() {
	stripe.Key = resource.Resource.StripePublic.Value
}

type FingerprintRequest struct {
	Fingerprint string `json:"fingerprint"`
}

type UserCredentials struct {
	UserID      string `json:"userID"`
	AccessToken string `json:"accessToken"`
}

func GetErrorMessage(err error) string {
	if apiError, ok := err.(*terminal.Error); ok {
		return strings.Trim(apiError.JSON.ExtraFields["message"].Raw(), "\"")
	} else {
		return err.Error()
	}
}

func (u UserCredentials) String() string {
	return fmt.Sprintf("{UserID: '%s', AccessToken: '%s...'}", u.UserID, string(u.AccessToken[:10]))
}

func FetchUserToken(publicKey string) (*UserCredentials, error) {
	fingerprint := FingerprintRequest{Fingerprint: publicKey}

	resp, body, err := authPost("ssh/login", fingerprint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 500 {
		return nil, errors.New(fmt.Sprintf("Server error: %s", string(body)))
	}

	var creds UserCredentials
	if err := json.Unmarshal(body, &creds); err != nil {
		return nil, errors.New(fmt.Sprintf("Error: %s, Body: %s", err, body))
	}

	if creds.AccessToken == "" {
		return nil, errors.New(fmt.Sprintf("Failed to fetch: %s", publicKey))
	}

	return &creds, nil
}

func StripeCreditCard(card *stripe.CardParams) (*stripe.Token, *string) {
	tokenParams := &stripe.TokenParams{Card: card}
	tokenResult, err := token.New(tokenParams)

	if err != nil {
		error := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			error = stripeErr.Msg
		} else {
			error = err.Error()
		}
		return tokenResult, &error
	}

	return tokenResult, nil
}

func authPost(path string, payload any) (*http.Response, []byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	url := fmt.Sprintf("%s/%s", resource.Resource.AuthWorker.Url, path)
	request, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+resource.Resource.AuthFingerprintKey.Value)
	resp, err := http.DefaultClient.Do(request)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}
