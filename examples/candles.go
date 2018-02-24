package main

import (
	"fmt"
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)
	history := oanda.GetCandles("EUR_USD")
	fmt.Println(history)
}
