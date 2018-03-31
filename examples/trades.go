package main

import (
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
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
	spew.Dump("%+v\n", tradesResponse)

	openTrades := oanda.GetOpenTrades()
	spew.Dump("%+v\n", openTrades)

	trade := oanda.GetTradesForInstrument("AUD_USD")
	spew.Dump("%+v\n", trade)
}
