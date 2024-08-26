package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetOpenPositions(t *testing.T) {
	defer logTestResult(t, "GetOpenPositions")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/openPositions" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := OpenPositions{
			LastTransactionID: "1000",
			Positions: []struct {
				Instrument string `json:"instrument"`
				Long       struct {
					AveragePrice string   `json:"averagePrice"`
					Pl           string   `json:"pl"`
					ResettablePL string   `json:"resettablePL"`
					TradeIDs     []string `json:"tradeIDs"`
					Units        string   `json:"units"`
					UnrealizedPL string   `json:"unrealizedPL"`
				} `json:"long"`
				Pl           string `json:"pl"`
				ResettablePL string `json:"resettablePL"`
				Short        struct {
					AveragePrice string   `json:"averagePrice"`
					Pl           string   `json:"pl"`
					ResettablePL string   `json:"resettablePL"`
					TradeIDs     []string `json:"tradeIDs"`
					Units        string   `json:"units"`
					UnrealizedPL string   `json:"unrealizedPL"`
				} `json:"short"`
				UnrealizedPL string `json:"unrealizedPL"`
			}{
				{
					Instrument: "EUR_USD",
					Long: struct {
						AveragePrice string   `json:"averagePrice"`
						Pl           string   `json:"pl"`
						ResettablePL string   `json:"resettablePL"`
						TradeIDs     []string `json:"tradeIDs"`
						Units        string   `json:"units"`
						UnrealizedPL string   `json:"unrealizedPL"`
					}{
						AveragePrice: "1.1000",
						Pl:           "10.00",
						ResettablePL: "10.00",
						TradeIDs:     []string{"1", "2"},
						Units:        "100",
						UnrealizedPL: "5.00",
					},
					Pl:           "10.00",
					ResettablePL: "10.00",
					UnrealizedPL: "5.00",
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

	positions, err := c.GetOpenPositions()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if positions.LastTransactionID != "1000" {
		t.Errorf("Expected LastTransactionID to be 1000, got %s", positions.LastTransactionID)
	}

	if len(positions.Positions) != 1 {
		t.Fatalf("Expected 1 position, got %d", len(positions.Positions))
	}

	position := positions.Positions[0]
	if position.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", position.Instrument)
	}
	if position.Long.AveragePrice != "1.1000" {
		t.Errorf("Expected Long.AveragePrice to be 1.1000, got %s", position.Long.AveragePrice)
	}
	if position.Long.Units != "100" {
		t.Errorf("Expected Long.Units to be 100, got %s", position.Long.Units)
	}
}

func TestClosePosition(t *testing.T) {
	defer logTestResult(t, "ClosePosition")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/positions/EUR_USD/close" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		var payload ClosePositionPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response := ModifiedTrade{
			LastTransactionID: "1000",
			OrderCreateTransaction: struct {
				Type         string    `json:"type"`
				Instrument   string    `json:"instrument"`
				Units        string    `json:"units"`
				TimeInForce  string    `json:"timeInForce"`
				PositionFill string    `json:"positionFill"`
				Reason       string    `json:"reason"`
				TradeClose   struct {
					Units   string `json:"units"`
					TradeID string `json:"tradeID"`
				} `json:"tradeClose"`
				ID        string    `json:"id"`
				UserID    int       `json:"userID"`
				AccountID string    `json:"accountID"`
				BatchID   string    `json:"batchID"`
				RequestID string    `json:"requestID"`
				Time      time.Time `json:"time"`
			}{
				Type:       "MARKET_ORDER",
				Instrument: "EUR_USD",
				Units:      "-100",
				Reason:     "POSITION_CLOSEOUT",
				ID:         "1000",
				AccountID:  "test-account",
				Time:       time.Now(),
			},
			OrderFillTransaction: struct {
				Type           string    `json:"type"`
				Instrument     string    `json:"instrument"`
				Units          string    `json:"units"`
				Price          string    `json:"price"`
				FullPrice      FullPrice `json:"fullPrice"`
				PL             string    `json:"pl"`
				Financing      string    `json:"financing"`
				Commission     string    `json:"commission"`
				AccountBalance string    `json:"accountBalance"`
				TradeOpened    string    `json:"tradeOpened"`
				TimeInForce    string    `json:"timeInForce"`
				PositionFill   string    `json:"positionFill"`
				Reason         string    `json:"reason"`
				TradesClosed   []struct {
					TradeID    string `json:"tradeID"`
					Units      string `json:"units"`
					RealizedPL string `json:"realizedPL"`
					Financing  string `json:"financing"`
				} `json:"tradesClosed"`
				TradeReduced struct {
					TradeID    string `json:"tradeID"`
					Units      string `json:"units"`
					RealizedPL string `json:"realizedPL"`
					Financing  string `json:"financing"`
				} `json:"tradeReduced"`
				ID            string    `json:"id"`
				UserID        int       `json:"userID"`
				AccountID     string    `json:"accountID"`
				BatchID       string    `json:"batchID"`
				RequestID     string    `json:"requestID"`
				OrderID       string    `json:"orderId"`
				ClientOrderID string    `json:"clientOrderId"`
				Time          time.Time `json:"time"`
			}{
				Type:       "ORDER_FILL",
				Instrument: "EUR_USD",
				Units:      "-100",
				Price:      "1.1000",
				PL:         "10.00",
				ID:         "1001",
				AccountID:  "test-account",
				Time:       time.Now(),
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

	payload := ClosePositionPayload{
		LongUnits:  "100",
		ShortUnits: "0",
	}

	modifiedTrade, err := c.ClosePosition("EUR_USD", payload)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if modifiedTrade.LastTransactionID != "1000" {
		t.Errorf("Expected LastTransactionID to be 1000, got %s", modifiedTrade.LastTransactionID)
	}

	if modifiedTrade.OrderCreateTransaction.Type != "MARKET_ORDER" {
		t.Errorf("Expected OrderCreateTransaction.Type to be MARKET_ORDER, got %s", modifiedTrade.OrderCreateTransaction.Type)
	}

	if modifiedTrade.OrderCreateTransaction.Instrument != "EUR_USD" {
		t.Errorf("Expected OrderCreateTransaction.Instrument to be EUR_USD, got %s", modifiedTrade.OrderCreateTransaction.Instrument)
	}

	if modifiedTrade.OrderCreateTransaction.Units != "-100" {
		t.Errorf("Expected OrderCreateTransaction.Units to be -100, got %s", modifiedTrade.OrderCreateTransaction.Units)
	}

	if modifiedTrade.OrderFillTransaction.Type != "ORDER_FILL" {
		t.Errorf("Expected OrderFillTransaction.Type to be ORDER_FILL, got %s", modifiedTrade.OrderFillTransaction.Type)
	}

	if modifiedTrade.OrderFillTransaction.Price != "1.1000" {
		t.Errorf("Expected OrderFillTransaction.Price to be 1.1000, got %s", modifiedTrade.OrderFillTransaction.Price)
	}
}
