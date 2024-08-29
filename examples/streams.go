package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/awoldes/goanda"
	"github.com/joho/godotenv"
)

func streams() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")

	// Create a new connection
	config := &goanda.ConnectionConfig{
		Live: os.Getenv("OANDA_LIVE") == "false",
	}
	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	// Create a new streaming connection
	streaming := oanda.NewStreamingConnection()

	// Create a WaitGroup to manage our goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Create a channel to signal goroutines to stop
	stop := make(chan struct{})

	// Start price streaming
	go func() {
		defer wg.Done()
		instruments := []string{"EUR_USD", "USD_JPY", "GBP_USD"}
		err := streaming.StreamPrices(instruments, func(response goanda.PricingStreamResponse) {
			fmt.Printf("Price update: %s - %s - Bid: %s, Ask: %s\n",
				response.Time,
				response.Instrument,
				response.Bids[0].Price,
				response.Asks[0].Price)
		})
		if err != nil {
			log.Printf("Error streaming prices: %v", err)
		}
	}()

	// Start transaction streaming
	go func() {
		defer wg.Done()
		err := streaming.StreamTransactions(func(response goanda.TransactionStreamResponse) {
			fmt.Printf("Transaction update: %s - Type: %s, ID: %s\n",
				response.Time,
				response.Type,
				response.TransactionID)
		})
		if err != nil {
			log.Printf("Error streaming transactions: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal. Shutting down...")

	// Signal goroutines to stop
	close(stop)

	// Wait for goroutines to finish
	wg.Wait()

	fmt.Println("Streaming example completed.")
}
