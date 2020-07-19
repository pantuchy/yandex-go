package yandex

import (
	"encoding/json"
	"log"
)

type Store struct {
	AccountId            string   `json:"account_id"`
	Test                 bool     `json:"test"`
	FiscalizationEnabled bool     `json:"fiscalization_enabled"`
	PaymentMethods       []string `json:"payment_methods"`
}

func (y *Yandex) GetStoreInfo() (*Store, error) {
	r := &HttpRequest{
		Method:     "GET",
		Path:       "/me",
		OAuthToken: y.OAuthToken,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Store{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}
