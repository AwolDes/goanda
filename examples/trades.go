package main

import (
	"fmt"
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/joho/godotenv"
)

func getTrades() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	tradesResponse := oanda.GetTrade("54")
	fmt.Printf("%+v\n", tradesResponse)

	openTrades := oanda.GetOpenTrades()
	fmt.Printf("%+v\n", openTrades)

	trade := oanda.GetTradesForInstrument("AUD_USD")
	fmt.Printf("%+v\n", trade)
}
