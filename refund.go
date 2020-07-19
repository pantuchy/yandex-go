package yandex

import (
	"encoding/json"
	"log"

	"github.com/google/go-querystring/query"
)

type Source struct {
	AccountId string  `json:"account_id"`
	Amount    *Amount `json:"amount"`
}

type Refund struct {
	Id          string    `json:"id"`
	PaymentId   string    `json:"payment_id"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	Amount      *Amount   `json:"amount"`
	Description string    `json:"description,omitempty"`
	Sources     []*Source `json:"sources,omitempty"`
}

type RefundRequest struct {
	PaymentId   string    `json:"payment_id"`
	Amount      *Amount   `json:"amount"`
	Description string    `json:"description,omitempty"`
	Receipt     *Receipt  `json:"receipt,omitempty"`
	Sources     []*Source `json:"sources,omitempty"`
}

func (y *Yandex) CreateRefund(idempKey string, req *RefundRequest) (*Refund, error) {
	q, err := query.Values(req)

	if err != nil {
		log.Printf("Failed creating query: %v\n", err)
		return nil, err
	}

	r := &HttpRequest{
		Method:         "POST",
		Path:           "/refunds",
		ShopId:         y.ShopId,
		SecretKey:      y.SecretKey,
		IdempotenceKey: idempKey,
		Data:           q,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Refund{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) GetRefundInfo(id string) (*Refund, error) {
	r := &HttpRequest{
		Method:    "GET",
		Path:      "/refunds/" + id,
		ShopId:    y.ShopId,
		SecretKey: y.SecretKey,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Refund{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}
