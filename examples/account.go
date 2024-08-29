package main

import (
	"log"
	"os"
	"time"
	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getAccountDetails() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &goanda.ConnectionConfig{
		UserAgent: "goanda",
		Timeout: 10 * time.Second,
		Live: false,
	}

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	accountChanges, err := oanda.GetAccountChanges("101-011-6559702-001", "54")
	spew.Dump(accountChanges)

	accountInstruments, err := oanda.GetAccountInstruments("101-011-6559702-001")
	spew.Dump(accountInstruments)

	accountSummary, err := oanda.GetAccountSummary()
	spew.Dump(accountSummary)

	account, err := oanda.GetAccount("101-011-6559702-003")
	spew.Dump(account)

	accounts, err := oanda.Accounts()
	spew.Dump(accounts)
}
