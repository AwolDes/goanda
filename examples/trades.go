package main

import (
	"log"
	"os"
	"time"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getTrades() {
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

	tradesResponse, err := oanda.GetTrade("54")
	if err != nil {
		log.Fatalf("Error getting trade: %v", err)
	}
	spew.Dump("%+v\n", tradesResponse)

	openTrades, err := oanda.GetOpenTrades()
	if err != nil {
		log.Fatalf("Error getting open trades: %v", err)
	}
	spew.Dump("%+v\n", openTrades)

	trade, err := oanda.GetTradesForInstrument("AUD_USD")
	if err != nil {
		log.Fatalf("Error getting trades for instrument: %v", err)
	}
	spew.Dump("%+v\n", trade)

	reduceTrade, err:= oanda.ReduceTradeSize("AUD_USD", goanda.CloseTradePayload{
		Units: "100.00",
	})
	if err != nil {
		log.Fatalf("Error reducing trade size: %v", err)
	}

	spew.Dump("%+v\n", reduceTrade)
}
