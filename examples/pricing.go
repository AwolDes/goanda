package main

import (
	"log"
	"os"

	"github.com/kuroko1t/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getPricing() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	instruments := []string{"AUD_USD", "EUR_NZD"}
	orderResponse := oanda.GetPricingForInstruments(instruments)
	spew.Dump("%+v\n", orderResponse)
}
