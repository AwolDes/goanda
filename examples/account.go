package main

import (
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func getAccountDetails() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	accountChanges := oanda.GetAccountChanges("101-011-6559702-001", "54")
	spew.Dump(accountChanges)

	accountInstruments := oanda.GetAccountInstruments("101-011-6559702-001")
	spew.Dump(accountInstruments)

	accountSummary := oanda.GetAccountSummary("101-011-6559702-001")
	spew.Dump(accountSummary)

	account := oanda.GetAccount("101-011-6559702-003")
	spew.Dump(account)

	accounts := oanda.GetAccounts()
	spew.Dump(accounts)
}
