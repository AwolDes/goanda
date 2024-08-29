package main

import (
	"log"
	"os"
	"time"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func placeMarketOrder() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &goanda.ConnectionConfig{
		UserAgent: "goanda",
		Timeout:   10 * time.Second,
		Live:      false,
	}

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")

	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	order := goanda.OrderPayload{
		Order: goanda.OrderBody{
			Units:        10000,
			Instrument:   "EUR_USD",
			TimeInForce:  "GTC",
			Type:         "MARKET",
			PositionFill: "DEFAULT",
			Price:        "1.25000",
		},
	}

	orderResult, err := oanda.CreateOrder(order)
	if err != nil {
		log.Fatalf("Error creating order: %v", err)
	}

	spew.Dump("%+v\n", orderResult)
}
