package resource

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type resource struct {
	OpenApiWorker struct {
		Url string `json:"url"`
	}
	AuthWorker struct {
		Url string `json:"url"`
	}
	StripePublic struct {
		Value string `json:"value"`
	}
	AuthFingerprintKey struct {
		Value string `json:"value"`
	}
	SSHKey struct {
		Public  string `json:"public"`
		Private string `json:"private"`
	}
}

var Resource resource

func init() {
	val := reflect.ValueOf(&Resource).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		envVarName := fmt.Sprintf("SST_RESOURCE_%s", typeField.Name)
		envValue, exists := os.LookupEnv(envVarName)
		if !exists {
			panic(fmt.Sprintf("Environment variable %s is required", envVarName))
		}
		if err := json.Unmarshal([]byte(envValue), field.Addr().Interface()); err != nil {
			panic(err)
		}
	}
}
