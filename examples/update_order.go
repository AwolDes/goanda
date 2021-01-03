package main

import (
	"log"
	"os"

	"github.com/kuroko1t/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func updateOrder() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)
	order := goanda.OrderPayload{
		Order: goanda.OrderBody{
			Units:        "1000",
			Instrument:   "EUR_USD",
			TimeInForce:  "FOK",
			Type:         "LIMIT",
			PositionFill: "DEFAULT",
			Price:        "1.25000",
			StopLossOnFill: &goanda.OnFill{
				TimeInForce: "GTC",
				Price:       "1.2000",
			},
		},
	}

	orderResponse := oanda.UpdateOrder("order-id", order)
	spew.Dump("%+v\n", orderResponse)
}
