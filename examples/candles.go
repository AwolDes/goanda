package main

import (
	"log"
	"os"
	"time"
	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func candles() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &goanda.ConnectionConfig{
		UserAgent: "goanda",
		Timeout: 10 * time.Second,
		Live: false,
	}

	granularity := goanda.GranularityFiveSeconds

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")

	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	history, err := oanda.GetCandles("EUR_USD", 10, granularity)
	if err != nil {
		log.Fatalf("Error getting candles: %v", err)
	}

	spew.Dump(history)
}
