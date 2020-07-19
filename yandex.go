package yandex

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
)

type Yandex struct {
	ShopId    string
	SecretKey string
}

type HttpRequest struct {
	Path           string
	Method         string
	ShopId         string
	SecretKey      string
	IdempotenceKey string
	Data           url.Values
}

type ErrorResponse struct {
	Type        string `json:"type"`
	Id          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Parameter   string `json:"parameter"`
}

type Amount struct {
	Value    decimal.Decimal `json:"value"`
	Currency string          `json:"currency"`
}

func (r *HttpRequest) SendRequest() ([]byte, error) {
	const baseURL string = "https://payment.yandex.net/api/v3"

	c := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(r.ShopId+":"+r.SecretKey)))
	req.Header.Set("Idempotence-Key", r.IdempotenceKey)
	req.Header.SetRequestURI(baseURL + r.Path)
	req.Header.SetMethod(strings.ToUpper(r.Method))
	req.Header.SetUserAgent("Mozilla/4.0 (compatible; Golang Yandex API)")
	req.Header.SetContentType("application/json")

	if r.Method == "GET" {
		req.SetRequestURI(fmt.Sprintf("%s%s?%s", baseURL, r.Path, r.Data.Encode()))
	} else {
		req.SetBody([]byte(r.Data.Encode()))
	}

	if err := c.Do(req, res); err != nil {
		return nil, &Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	if res.StatusCode() != 200 {
		e := &ErrorResponse{}

		if err := json.Unmarshal(res.Body(), e); err != nil {
			log.Printf("Failed unmarshaling bytes to struct: %v\n", err)
			return nil, err
		}

		return nil, &Error{
			Code:    res.StatusCode(),
			ApiCode: e.Code,
			Message: e.Description,
		}
	}

	return res.Body(), nil
}
