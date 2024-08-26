package main

import (
	"log"
	"os"
	"time"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getOrders() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &goanda.ConnectionConfig{
		UserAgent: "goanda",
		Timeout: 10 * time.Second,
		Live: false,
	}

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")

	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	orders, err := oanda.GetOrders("EUR_USD")
	if err != nil {
		log.Fatalf("Error getting orders: %v", err)
	}
	
	pendingOrders, err := oanda.GetPendingOrders()
	if err != nil {
		log.Fatalf("Error getting pending orders: %v", err)
	}
	spew.Dump("%+v\n", orders)
	spew.Dump("%+v\n", pendingOrders)
}
