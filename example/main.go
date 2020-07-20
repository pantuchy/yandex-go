package main

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/pantuchy/yandex-go"
	"github.com/shopspring/decimal"
)

func main() {
	price, err := decimal.NewFromString("120.00")

	if err != nil {
		log.Fatalf("Failed parsing string to decimal: %v\n", err)
	}

	req := &yandex.PaymentRequest{
		Amount: &yandex.Amount{
			Value:    price,
			Currency: "EUR",
		},
		Confirmation: &yandex.Confirmation{
			Type:      "redirect",
			ReturnUrl: "https://www.merchant-website.com/return_url",
		},
	}

	card := &yandex.Card{
		Number:      "1234567812345678",
		ExpiryYear:  "2022",
		ExpiryMonth: "07",
		CSC:         "123",
		CardHolder:  "John Travola",
	}

	client := &yandex.Yandex{
		ShopId:     "123456", // Required for payments, refunds and receipts
		SecretKey:  "qwerty", // Required for payments, refunds and receipts
		OAuthToken: "token",  // Required only for Webhooks and Store information
	}

	idempKey, err := uuid.NewV4()

	if err != nil {
		log.Fatalf("Failed to generate Idempotence-Key: %v\n", err)
	}

	res, err := client.CreatePayment(idempKey.String(), req.WithBankCard(card))

	if err != nil {
		log.Fatalf("Creating payment failed: %v\n", err)
	}

	fmt.Println(res)
}
