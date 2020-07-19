package yandex

import (
	"encoding/json"
	"log"

	"github.com/google/go-querystring/query"
)

type Webhook struct {
	Id    string `json:"id"`
	Event string `json:"event"`
	Url   string `json:"url"`
}

type WebhooksListResponse struct {
	Type  string     `json:"type"`
	Items []*Webhook `json:"items"`
}

func (y *Yandex) SubscribeToWebhook(idempKey string, req *Webhook) (*Webhook, error) {
	q, err := query.Values(req)

	if err != nil {
		log.Printf("Failed creating query: %v\n", err)
		return nil, err
	}

	r := &HttpRequest{
		Method:         "POST",
		Path:           "/webhooks",
		IdempotenceKey: idempKey,
		OAuthToken:     y.OAuthToken,
		Data:           q,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Webhook{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) GetWebhooksList() (*WebhooksListResponse, error) {
	r := &HttpRequest{
		Method:     "GET",
		Path:       "/webhooks",
		OAuthToken: y.OAuthToken,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &WebhooksListResponse{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) DeleteWebhook(id string) error {
	r := &HttpRequest{
		Method:     "DELETE",
		Path:       "/webhooks/" + id,
		OAuthToken: y.OAuthToken,
	}

	_, err := r.SendRequest()

	if err != nil {
		return err
	}

	return nil
}
