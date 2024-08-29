package main

import (
	"log"
	"os"
	"time"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func positions() {
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

	closePosition, err := oanda.ClosePosition("AUD_USD", goanda.ClosePositionPayload{
		LongUnits:  "NONE",
		ShortUnits: "ALL",
	})
	if err != nil {
		log.Fatalf("Error closing position: %v", err)
	}

	spew.Dump("%+v\n", closePosition)

	openPositions, err := oanda.GetOpenPositions()
	if err != nil {
		log.Fatalf("Error getting open positions: %v", err)
	}

	spew.Dump("%+v\n", openPositions)
}
