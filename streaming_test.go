package goanda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestNewStreamingConnection(t *testing.T) {
	defer logTestResult(t, "TestNewStreamingConnection")
	conn := &Connection{
		hostname: "https://api-fxpractice.oanda.com/v3",
	}
	sc := NewStreamingConnection(conn)

	if sc.streamURL != "https://stream-fxpractice.oanda.com/v3" {
		t.Errorf("Expected streamURL to be https://stream-fxpractice.oanda.com/v3, got %s", sc.streamURL)
	}

	conn.hostname = "https://api-fxtrade.oanda.com/v3"
	sc = NewStreamingConnection(conn)

	if sc.streamURL != "https://stream-fxtrade.oanda.com/v3" {
		t.Errorf("Expected streamURL to be https://stream-fxtrade.oanda.com/v3, got %s", sc.streamURL)
	}
}

func TestStreamPrices(t *testing.T) {
	defer logTestResult(t, "TestStreamPrices")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/pricing/stream" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		instruments := r.URL.Query().Get("instruments")
		if instruments != "EUR_USD" {
			t.Errorf("Unexpected instruments: %s", instruments)
			http.Error(w, "Invalid instruments", http.StatusBadRequest)
			return
		}

		response := PricingStreamResponse{
			Type:       "PRICE",
			Time:       time.Now().Format(time.RFC3339),
			Instrument: "EUR_USD",
			Bids: []struct {
				Price     string `json:"price"`
				Liquidity int    `json:"liquidity"`
			}{{Price: "1.1000", Liquidity: 1000000}},
			Asks: []struct {
				Price     string `json:"price"`
				Liquidity int    `json:"liquidity"`
			}{{Price: "1.1001", Liquidity: 1000000}},
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	conn := &Connection{
		hostname:   server.URL,
		accountID:  "test-account",
		authHeader: "Bearer test-token",
		client:     *server.Client(),
	}
	sc := NewStreamingConnection(conn)

	// Override the streamURL to use the test server
	sc.streamURL = server.URL

	instruments := []string{"EUR_USD"}
	err := sc.StreamPrices(instruments, func(response PricingStreamResponse) {
		if response.Type != "PRICE" {
			t.Errorf("Expected response type to be PRICE, got %s", response.Type)
		}
		if response.Instrument != "EUR_USD" {
			t.Errorf("Expected instrument to be EUR_USD, got %s", response.Instrument)
		}
		if len(response.Bids) == 0 {
			t.Errorf("Expected non-empty Bids slice")
		} else if response.Bids[0].Price != "1.1000" {
			t.Errorf("Expected bid price to be 1.1000, got %s", response.Bids[0].Price)
		}
		if len(response.Asks) == 0 {
			t.Errorf("Expected non-empty Asks slice")
		} else if response.Asks[0].Price != "1.1001" {
			t.Errorf("Expected ask price to be 1.1001, got %s", response.Asks[0].Price)
		}
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestStreamTransactions(t *testing.T) {
	defer logTestResult(t, "TestStreamTransactions")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/transactions/stream" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := TransactionStreamResponse{
			Type:          "TRANSACTION",
			Time:          time.Now().Format(time.RFC3339),
			TransactionID: "1234",
			AccountID:     "test-account",
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	conn := &Connection{
		hostname:   server.URL,
		accountID:  "test-account",
		authHeader: "Bearer test-token",
		client:     *server.Client(),
	}
	sc := NewStreamingConnection(conn)

	// Override the streamURL to use the test server
	sc.streamURL = server.URL

	err := sc.StreamTransactions(func(response TransactionStreamResponse) {
		if response.Type != "TRANSACTION" {
			t.Errorf("Expected response type to be TRANSACTION, got %s", response.Type)
		}
		if response.AccountID != "test-account" {
			t.Errorf("Expected account ID to be test-account, got %s", response.AccountID)
		}
		if response.TransactionID != "1234" {
			t.Errorf("Expected transaction ID to be 1234, got %s", response.TransactionID)
		}
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestStreamHeartbeat(t *testing.T) {
	defer logTestResult(t, "TestStreamHeartbeat")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HeartbeatResponse{
			Type: "HEARTBEAT",
			Time: time.Now().Format(time.RFC3339),
		}
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	conn := &Connection{
		hostname:   server.URL,
		accountID:  "test-account",
		authHeader: "Bearer test-token",
		client:     *server.Client(),
	}
	sc := NewStreamingConnection(conn)

	// Override the streamURL to use the test server
	sc.streamURL = server.URL

	// This test is a bit tricky because heartbeats are handled internally.
	// We'll use the StreamPrices function, but send a heartbeat instead.
	err := sc.StreamPrices([]string{"EUR_USD"}, func(response PricingStreamResponse) {
		t.Errorf("Unexpected pricing response: %+v", response)
	})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// If we reach this point without errors, it means the heartbeat was properly handled
}

func TestStreamingIntegration(t *testing.T) {
	defer logTestResult(t, "TestStreamingIntegration")
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	if apiKey == "" || accountID == "" {
		t.Fatal("OANDA_API_KEY and OANDA_ACCOUNT_ID must be set in .env file")
	}

	// Create a new connection
	config := &ConnectionConfig{
		UserAgent: "goanda-test",
		Timeout:   10 * time.Second,
		Live:      false, // Use demo account
	}
	conn, err := NewConnection(accountID, apiKey, config)
	if err != nil {
		t.Fatalf("Error creating connection: %v", err)
	}

	sc := conn.NewStreamingConnection()

	// Test StreamPrices
	t.Run("StreamPrices", func(t *testing.T) {
		defer logTestResult(t, "TestStreamPrices")
		instruments := []string{"EUR_USD", "GBP_USD"}
		receivedPrice := make(chan bool, 1)
		errorChan := make(chan error, 1)

		go func() {
			err := sc.StreamPrices(instruments, func(response PricingStreamResponse) {
				if response.Type == "PRICE" {
					select {
					case receivedPrice <- true:
					default:
					}
				}
			})
			if err != nil {
				errorChan <- fmt.Errorf("Error streaming prices: %v", err)
			}
		}()

		select {
		case <-receivedPrice:
			// Test passed
		case err := <-errorChan:
			t.Errorf("StreamPrices error: %v", err)
		case <-time.After(45 * time.Second):
			t.Error("Timeout waiting for price stream")
		}
	})

	time.Sleep(1 * time.Second)

	// Test StreamTransactions
	t.Run("StreamTransactions", func(t *testing.T) {
		defer logTestResult(t, "TestStreamTransactions")
		receivedTransaction := make(chan bool, 1)
		errorChan := make(chan error, 1)

		go func() {
			err := sc.StreamTransactions(func(response TransactionStreamResponse) {
				if response.Type == "ORDER_FILL" {
					receivedTransaction <- true
				}
			})
			if err != nil {
				errorChan <- fmt.Errorf("Error streaming transactions: %v", err)
			}
		}()
		time.Sleep(3 * time.Second)
		// Create a market order to trigger a transaction
		orderBody := OrderBody{
			Instrument:   "EUR_USD",
			Units:        100,
			Type:         "MARKET",
			TimeInForce:  "FOK",
			PositionFill: "DEFAULT",
		}
		order, err := conn.CreateOrder(OrderPayload{Order: orderBody})
		if err != nil {
			t.Fatalf("Error creating order: %v", err)
		}

		select {
		case <-receivedTransaction:
			// Test passed
		case err := <-errorChan:
			t.Errorf("StreamTransactions error: %v", err)
		case <-time.After(45 * time.Second):
			t.Error("Timeout waiting for transaction stream")
		}
		// Close the order
		_, err = conn.ReduceTradeSize(order.OrderFillTransaction.TradeOpened.TradeID, CloseTradePayload{
			Units: "ALL",
		})
		if err != nil {
			t.Fatalf("Error closing order: %v", err)
		}
	})
}
