package main

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/pantuchy/yandex-go"
)

func main() {
	req := &yandex.PaymentRequest{}
	client := &yandex.Yandex{}

	idempKey, err := uuid.NewV4()

	if err != nil {
		log.Fatalf("Failed to generate Idempotence-Key: %v\n", err)
	}

	res, err := client.CreatePayment(idempKey.String(), req.WithYandexMoney())

	if err != nil {
		log.Fatalf("Creating payment failed: %v\n", err)
	}

	fmt.Println(res)
}
