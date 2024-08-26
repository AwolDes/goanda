package main

import (
	"log"
	"os"
	"time"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getTransactions() {
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

	toTime := time.Now()
	// AddDate(Y, M, D) - https://golang.org/pkg/time/#Time.AddDate
	fromTime := toTime.AddDate(0, -1, 0)
	transactions, err := oanda.GetTransactions(fromTime, toTime)
	if err != nil {
		log.Fatalf("Error getting transactions: %v", err)
	}
	spew.Dump("%+v\n", transactions)

	transactionsSince, err := oanda.GetTransactionsSinceId("55")
	if err != nil {
		log.Fatalf("Error getting transactions since id: %v", err)
	}
	spew.Dump("%+v\n", transactionsSince)
}
