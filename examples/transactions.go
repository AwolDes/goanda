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
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	toTime := time.Now()
	// AddDate(Y, M, D) - https://golang.org/pkg/time/#Time.AddDate
	fromTime := toTime.AddDate(0, -1, 0)
	transactions := oanda.GetTransactions(fromTime, toTime)
	spew.Dump("%+v\n", transactions)

	transactionsSince := oanda.GetTransactionsSinceId("55")
	spew.Dump("%+v\n", transactionsSince)
}
