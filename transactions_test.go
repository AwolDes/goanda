package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	defer logTestResult(t, "GetTransactions")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/transactions" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		if from == "" || to == "" {
			t.Errorf("Missing from or to query parameters")
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		response := TransactionPages{
			From:              time.Now().Add(-24 * time.Hour),
			To:                time.Now(),
			PageSize:          100,
			Count:             2,
			Pages:             []string{"0", "1"},
			LastTransactionID: "1000",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "test-account",
		client:    *server.Client(),
	}

	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()

	transactions, err := c.GetTransactions(from, to)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if transactions.Count != 2 {
		t.Errorf("Expected Count to be 2, got %d", transactions.Count)
	}

	if transactions.LastTransactionID != "1000" {
		t.Errorf("Expected LastTransactionID to be 1000, got %s", transactions.LastTransactionID)
	}
}

func TestGetTransaction(t *testing.T) {
	defer logTestResult(t, "GetTransaction")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/transactions/1000" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := Transaction{
			LastTransactionID: "1000",
			Transaction: struct {
				AccountBalance string    `json:"accountBalance"`
				AccountID      string    `json:"accountID"`
				BatchID        string    `json:"batchID"`
				Financing      string    `json:"financing"`
				ID             string    `json:"id"`
				Instrument     string    `json:"instrument"`
				OrderID        string    `json:"orderID"`
				Pl             string    `json:"pl"`
				Price          string    `json:"price"`
				Reason         string    `json:"reason"`
				Time           time.Time `json:"time"`
				TradeOpened    struct {
					TradeID string `json:"tradeID"`
					Units   string `json:"units"`
				} `json:"tradeOpened"`
				Type   string `json:"type"`
				Units  string `json:"units"`
				UserID int    `json:"userID"`
			}{
				ID:         "1000",
				Type:       "MARKET_ORDER",
				Instrument: "EUR_USD",
				Units:      "100",
				Price:      "1.1000",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "test-account",
		client:    *server.Client(),
	}

	transaction, err := c.GetTransaction("1000")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if transaction.LastTransactionID != "1000" {
		t.Errorf("Expected LastTransactionID to be 1000, got %s", transaction.LastTransactionID)
	}

	if transaction.Transaction.ID != "1000" {
		t.Errorf("Expected Transaction.ID to be 1000, got %s", transaction.Transaction.ID)
	}

	if transaction.Transaction.Type != "MARKET_ORDER" {
		t.Errorf("Expected Transaction.Type to be MARKET_ORDER, got %s", transaction.Transaction.Type)
	}
}

func TestGetTransactionsSinceId(t *testing.T) {
	defer logTestResult(t, "GetTransactionsSinceId")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/transactions/sinceid" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		id := r.URL.Query().Get("id")
		if id != "999" {
			t.Errorf("Unexpected id: %s", id)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		response := Transactions{
			LastTransactionID: "1001",
			Transactions: []struct {
				AccountBalance string    `json:"accountBalance"`
				AccountID      string    `json:"accountID"`
				BatchID        string    `json:"batchID"`
				Financing      string    `json:"financing"`
				ID             string    `json:"id"`
				Instrument     string    `json:"instrument"`
				OrderID        string    `json:"orderID"`
				Pl             string    `json:"pl"`
				Price          string    `json:"price"`
				Reason         string    `json:"reason"`
				Time           time.Time `json:"time"`
				TradeOpened    struct {
					TradeID string `json:"tradeID"`
					Units   string `json:"units"`
				} `json:"tradeOpened"`
				Type   string `json:"type"`
				Units  string `json:"units"`
				UserID int    `json:"userID"`
			}{
				{
					ID:         "1000",
					Type:       "MARKET_ORDER",
					Instrument: "EUR_USD",
					Units:      "100",
					Price:      "1.1000",
				},
				{
					ID:         "1001",
					Type:       "TRADE_CLOSE",
					Instrument: "EUR_USD",
					Units:      "-100",
					Price:      "1.1010",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "test-account",
		client:    *server.Client(),
	}

	transactions, err := c.GetTransactionsSinceId("999")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if transactions.LastTransactionID != "1001" {
		t.Errorf("Expected LastTransactionID to be 1001, got %s", transactions.LastTransactionID)
	}

	if len(transactions.Transactions) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(transactions.Transactions))
	}

	if transactions.Transactions[0].ID != "1000" {
		t.Errorf("Expected first transaction ID to be 1000, got %s", transactions.Transactions[0].ID)
	}

	if transactions.Transactions[1].ID != "1001" {
		t.Errorf("Expected second transaction ID to be 1001, got %s", transactions.Transactions[1].ID)
	}
}
