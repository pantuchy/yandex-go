package yandex

import (
	"encoding/json"
	"log"

	"github.com/google/go-querystring/query"
)

type VATData struct {
	Type   string  `json:"type"`
	Amount *Amount `json:"amount,omitempty"`
	Rate   string  `json:"rate,omitempty"`
}

type PayerBankDetails struct {
	FullName   string `json:"full_name"`
	ShortName  string `json:"short_name"`
	Address    string `json:"address"`
	INN        string `json:"inn"`
	KPP        string `json:"kpp"`
	BankName   string `json:"bank_name"`
	BankBranch string `json:"bank_branch"`
	BIC        string `json:"bank_bik"`
	Account    string `json:"account"`
}

type Card struct {
	Number         string `json:"number,omitempty"`
	ExpiryYear     string `json:"expiry_year"`
	ExpiryMonth    string `json:"expiry_month"`
	CSC            string `json:"csc,omitempty"`
	CardHolder     string `json:"cardholder,omitempty"`
	BIN            string `json:"first6,omitempty"`
	LastFourDigits string `json:"last4,omitempty"`
	CardType       string `json:"card_type,omitempty"`
	IssuerCountry  string `json:"issuer_country,omitempty"`
	IssuerName     string `json:"issuer_name,omitempty"`
	Source         string `json:"source,omitempty"`
}

type PaymentMethod struct {
	Type               string            `json:"type"`
	Id                 string            `json:"id,omitempty"`
	Saved              bool              `json:"saved,omitempty"`
	Title              string            `json:"title,omitempty"`
	Login              string            `json:"login,omitempty"`
	Phone              string            `json:"phone,omitempty"`
	Card               *Card             `json:"card,omitempty"`
	PaymentPurpose     string            `json:"payment_purpose,omitempty"`
	VATData            *VATData          `json:"vat_data,omitempty"`
	PaymentData        string            `json:"payment_data,omitempty"`
	PaymentMethodToken string            `json:"payment_method_token,omitempty"`
	PayerBankDetails   *PayerBankDetails `json:"payer_bank_details,omitempty"`
	AccountNumber      string            `json:"account_number,omitempty"`
}

type Leg struct {
	DepartureAirport   string `json:"departure_airport"`
	DestinationAirport string `json:"destination_airport"`
	DepartureDate      string `json:"departure_date"`
	CarrierCode        string `json:"carrier_code,omitempty"`
}

type Passenger struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Airline struct {
	TicketNumber     string       `json:"ticket_number,omitempty"`
	BookingReference string       `json:"booking_reference,omitempty"`
	Passengers       []*Passenger `json:"passengers,omitempty"`
	Legs             []*Leg       `json:"legs,omitempty"`
}

type Transfer struct {
	AccountId string  `json:"account_id"`
	Amount    *Amount `json:"amount"`
	Status    string  `json:"status,omitempty"`
}

type Confirmation struct {
	Type              string `json:"type"`
	Enforce           bool   `json:"enforce,omitempty"`
	Locale            string `json:"locale,omitempty"`
	ReturnUrl         string `json:"return_url,omitempty"`
	ConfirmationUrl   string `json:"confirmation_url,omitempty"`
	ConfirmationToken string `json:"confirmation_token,omitempty"`
	ConfirmationData  string `json:"confirmation_data,omitempty"`
}

type PaymentConfirmationRequest struct {
	Amount    *Amount     `json:"amount,omitempty"`
	Receipt   *Receipt    `json:"receipt,omitempty"`
	Airline   *Airline    `json:"airline,omitempty"`
	Transfers []*Transfer `json:"transfers,omitempty"`
}

type CancellationDetails struct {
	Party  string `json:"party"`
	Reason string `json:"reason"`
}

type AuthorizationDetails struct {
	RetrievalReferenceNumber string `json:"rrn,omitempty"`
	AuthCode                 string `json:"auth_code,omitempty"`
}

type Recipient struct {
	AccountId string `json:"account_id,omitempty"`
	GatewayId string `json:"gateway_id"`
}

type Payment struct {
	Id                   string                 `json:"id"`
	Status               string                 `json:"status"`
	Amount               *Amount                `json:"amount"`
	IncomeMmount         *Amount                `json:"income_amount,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Recipient            *Recipient             `json:"recipient"`
	PaymentMethod        *PaymentMethod         `json:"payment_method"`
	CapturedAt           string                 `json:"captured_at,omitempty"`
	CreatedAt            string                 `json:"created_at"`
	ExpiresAt            string                 `json:"expires_at,omitempty"`
	Confirmation         *Confirmation          `json:"confirmation,omitempty"`
	Test                 bool                   `json:"test"`
	RefundedAmount       *Amount                `json:"refunded_amount,omitempty"`
	Paid                 bool                   `json:"paid"`
	Refundable           bool                   `json:"refundable"`
	ReceiptRegistration  string                 `json:"receipt_registration,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	CancellationDetails  *CancellationDetails   `json:"cancellation_details,omitempty"`
	AuthorizationDetails *AuthorizationDetails  `json:"authorization_details,omitempty"`
	Transfers            []*Transfer            `json:"transfers,omitempty"`
}

type PaymentRequest struct {
	Amount            *Amount                `json:"amount"`
	Description       string                 `json:"description,omitempty"`
	Receipt           *Receipt               `json:"receipt,omitempty"`
	Recipient         *Recipient             `json:"recipient,omitempty"`
	PaymentToken      string                 `json:"payment_token,omitempty"`
	PaymentMethodId   string                 `json:"payment_method_id,omitempty"`
	PaymentMethodData *PaymentMethod         `json:"payment_method_data,omitempty"`
	Confirmation      *Confirmation          `json:"confirmation,omitempty"`
	SavePaymentMethod bool                   `json:"save_payment_method,omitempty"`
	Capture           bool                   `json:"capture,omitempty"`
	ClientIp          string                 `json:"client_ip,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	Airline           *Airline               `json:"airline,omitempty"`
	Transfers         []*Transfer            `json:"transfers,omitempty"`
}

func (y *Yandex) CreatePayment(idempKey string, req *PaymentRequest) (*Payment, error) {
	q, err := query.Values(req)

	if err != nil {
		log.Printf("Failed creating query: %v\n", err)
		return nil, err
	}

	r := &HttpRequest{
		Method:         "POST",
		Path:           "/payments",
		ShopId:         y.ShopId,
		SecretKey:      y.SecretKey,
		IdempotenceKey: idempKey,
		Data:           q,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Payment{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) GetPaymentInfo(id string) (*Payment, error) {
	r := &HttpRequest{
		Method:    "GET",
		Path:      "/payments/" + id,
		ShopId:    y.ShopId,
		SecretKey: y.SecretKey,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Payment{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) ConfirmPayment(idempKey, id string, req *PaymentConfirmationRequest) (*Payment, error) {
	q, err := query.Values(req)

	if err != nil {
		log.Printf("Failed creating query: %v\n", err)
		return nil, err
	}

	r := &HttpRequest{
		Method:         "POST",
		Path:           "/payments/" + id + "/capture",
		ShopId:         y.ShopId,
		SecretKey:      y.SecretKey,
		IdempotenceKey: idempKey,
		Data:           q,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Payment{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (y *Yandex) CancelPayment(idempKey, id string) (*Payment, error) {
	r := &HttpRequest{
		Method:         "POST",
		Path:           "/payments/" + id + "/cancel",
		ShopId:         y.ShopId,
		SecretKey:      y.SecretKey,
		IdempotenceKey: idempKey,
	}

	bytes, err := r.SendRequest()

	if err != nil {
		return nil, err
	}

	res := &Payment{}

	if err := json.Unmarshal(bytes, res); err != nil {
		log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (r *PaymentRequest) WithAlfaBank(login string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:  "alfabank",
		Login: login,
	}

	return r
}

func (r *PaymentRequest) WithPhoneBalance(phone string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:  "mobile_balance",
		Phone: phone,
	}

	return r
}

func (r *PaymentRequest) WithBankCard(card *Card) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "bank_card",
		Card: card,
	}

	return r
}

func (r *PaymentRequest) WithPartialPayment() *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "installments",
	}

	return r
}

func (r *PaymentRequest) WithCash(phone string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:  "cash",
		Phone: phone,
	}

	return r
}

func (r *PaymentRequest) WithSberbankB2B(purpose string, data *VATData) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:           "b2b_sberbank",
		PaymentPurpose: purpose,
		VATData:        data,
	}

	return r
}

func (r *PaymentRequest) WithSberbank(phone string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:  "sberbank",
		Phone: phone,
	}

	return r
}

func (r *PaymentRequest) WithTinkoffBank() *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "tinkoff_bank",
	}

	return r
}

func (r *PaymentRequest) WithYandexMoney() *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "yandex_money",
	}

	return r
}

func (r *PaymentRequest) WithApplePay(data string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:        "apple_pay",
		PaymentData: data,
	}

	return r
}

func (r *PaymentRequest) WithGooglePay(token string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:               "google_pay",
		PaymentMethodToken: token,
	}

	return r
}

func (r *PaymentRequest) WithQiwi(phone string) *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type:  "qiwi",
		Phone: phone,
	}

	return r
}

func (r *PaymentRequest) WithWeChat() *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "wechat",
	}

	return r
}

func (r *PaymentRequest) WithWebmoney() *PaymentRequest {
	r.PaymentMethodData = &PaymentMethod{
		Type: "webmoney",
	}

	return r
}
