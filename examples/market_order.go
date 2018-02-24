package main

import (
	"fmt"
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/joho/godotenv"
)

func placeMarketOrder() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)
	order := &goanda.Order{
		Order: goanda.OrderBody{
			Units:        "10000",
			Instrument:   "EUR_USD",
			TimeInForce:  "GTC",
			Type:         "MARKET",
			PositionFill: "DEFAULT",
			Price:        "1.25000",
		},
	}
	orderResult := oanda.CreateOrder(order)
	fmt.Printf("%+v\n", orderResult)
}
