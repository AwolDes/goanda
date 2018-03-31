package main

import (
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func positions() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	closePosition := oanda.ClosePosition("AUD_USD", goanda.ClosePositionPayload{
		LongUnits:  "NONE",
		ShortUnits: "ALL",
	})
	spew.Dump("%+v\n", closePosition)

	openPositions := oanda.GetOpenPositions()
	spew.Dump("%+v\n", openPositions)
}
