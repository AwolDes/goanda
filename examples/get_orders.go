package main

import (
	"log"
	"os"

	"github.com/kuroko1t/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getOrders() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)
	orders := oanda.GetOrders("EUR_USD")
	pendingOrders := oanda.GetPendingOrders()
	spew.Dump("%+v\n", orders)
	spew.Dump("%+v\n", pendingOrders)
}
