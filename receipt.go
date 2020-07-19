package yandex

import (
	"encoding/json"
	"log"

	"github.com/google/go-querystring/query"
	"github.com/shopspring/decimal"
)

type Customer struct {
	FullName string `json:"full_name,omitempty"`
	INN      string `json:"inn,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type Supplier struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	INN   string `json:"inn,omitempty"`
}

type Item struct {
	Description              string          `json:"description"`
	Quantity                 decimal.Decimal `json:"quantity"`
	Amount                   *Amount         `json:"amount"`
	VATCode                  uint32          `json:"vat_code"`
	PaymentSubject           string          `json:"payment_subject,omitempty"`
	PaymentMode              string          `json:"payment_mode,omitempty"`
	ProductCode              string          `json:"product_code,omitempty"`
	CountryOfOriginCode      string          `json:"country_of_origin_code,omitempty"`
	CustomsDeclarationNumber string          `json:"customs_declaration_number,omitempty"`
	Excise                   string          `json:"excise,omitempty"`
	Supplier                 *Supplier       `json:"supplier,omitempty"`
	AgentType                string          `json:"agent_type,omitempty"`
}

type Settlement struct {
	Type   string  `json:"type"`
	Amount *Amount `json:"amount"`
}

type Receipt struct {
	Id                   string        `json:"id,omitempty"`
	Type                 string        `json:"type,omitempty"`
	PaymentId            string        `json:"payment_id,omitempty"`
	RefundId             string        `json:"refund_id,omitempty"`
	Status               string        `json:"status,omitempty"`
	FiscalDocumentNumber string        `json:"fiscal_document_number,omitempty"`
	FiscalStorageNumber  string        `json:"fiscal_storage_number,omitempty"`
	FiscalAttribute      string        `json:"fiscal_attribute,omitempty"`
	RegisteredAt         string        `json:"registered_at,omitempty"`
	FiscalProviderId     string        `json:"fiscal_provider_id,omitempty"`
	Settlements          []*Settlement `json:"settlements,omitempty"`
	Customer             *Customer     `json:"customer,omitempty"`
	Items                []*Item       `json:"items"`
	TaxSystemCode        uint32        `json:"tax_system_code,omitempty"`
	Phone                string        `json:"phone,omitempty"`
	Email                string        `json:"email,omitempty"`
	OnBehalfOf           string        `json:"on_behalf_of,omitempty"`
}

type ReceiptRequest struct {
	Type          string        `json:"type"`
	PaymentId     string        `json:"payment_id,omitempty"`
	RefundId      string        `json:"refund_id,omitempty"`
	Customer      *Customer     `json:"customer"`
	Items         []*Item       `json:"items"`
	TaxSystemCode uint32        `json:"tax_system_code,omitempty"`
	Send          bool          `json:"send"`
	Settlements   []*Settlement `json:"settlements"`
	OnBehalfOf    string        `json:"on_behalf_of,omitempty"`
}

func (y *Yandex) CreateReceipt(idempKey string, req *ReceiptRequest) (*Receipt, error) {
	q, err := query.Values(req)

	if err != nil {
		log.Printf("Failed creating query: %v\n", err)
		return nil, err
	}

	r := &HttpRequest{
		Method:         "POST",
		Path:           "/receipts",
		ShopId:         y.ShopId,
		SecretKey:      y.SecretKey,
		IdempotenceKey: idempKey,
		Data:           q,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Receipt{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) GetReceiptInfo(id string) (*Receipt, error) {
	r := &HttpRequest{
		Method:    "GET",
		Path:      "/receipts/" + id,
		ShopId:    y.ShopId,
		SecretKey: y.SecretKey,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Receipt{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}
