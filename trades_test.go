package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTradesForInstrument(t *testing.T) {
	defer logTestResult(t, "TestGetTradesForInstrument")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/trades" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("instrument") != "EUR_USD" {
			t.Errorf("Unexpected instrument: %s", r.URL.Query().Get("instrument"))
			http.Error(w, "Invalid instrument", http.StatusBadRequest)
			return
		}

		response := ReceivedTrades{
			LastTransactionID: "1234",
			Trades: []Trade{
				{
					ID:                    "1",
					Instrument:            "EUR_USD",
					Price:                 "1.1000",
					OpenTime:              time.Now(),
					State:                 "OPEN",
					InitialUnits:          "100",
					InitialMarginRequired: "10.00",
					CurrentUnits:          "100",
					RealizedPL:            "0.00",
					UnrealizedPL:          "10.00",
					MarginUsed:            "10.00",
					AverageClosePrice:     "0.00",
					ClosingTransactionIDs: []string{},
					Financing:             "0.00",
					CloseTime:             time.Time{},
					ClientExtensions:      &OrderExtensions{Comment: "Test trade", Tag: "test"},
					TakeProfitOrder: &TakeProfitOrder{
						ID:               "2",
						CreateTime:       time.Now(),
						Type:             "TAKE_PROFIT",
						TradeID:          "1",
						Price:            "1.1100",
						TimeInForce:      "GTC",
						TriggerCondition: "DEFAULT",
						State:            "PENDING",
					},
					StopLossOrder: &StopLossOrder{
						ID:               "3",
						CreateTime:       time.Now(),
						Type:             "STOP_LOSS",
						TradeID:          "1",
						Price:            "1.0900",
						TimeInForce:      "GTC",
						TriggerCondition: "DEFAULT",
						State:            "PENDING",
						Guaranteed:       false,
					},
					TrailingStopLossOrder: &TrailingStopLossOrder{
						ID:               "4",
						CreateTime:       time.Now(),
						Type:             "TRAILING_STOP_LOSS",
						TradeID:          "1",
						Distance:         "0.0050",
						TimeInForce:      "GTC",
						TriggerCondition: "DEFAULT",
						State:            "PENDING",
					},
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

	trades, err := c.GetTradesForInstrument("EUR_USD")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if trades.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", trades.LastTransactionID)
	}

	if len(trades.Trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades.Trades))
	}

	trade := trades.Trades[0]
	if trade.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", trade.Instrument)
	}
	if trade.CurrentUnits != "100" {
		t.Errorf("Expected CurrentUnits to be 100, got %s", trade.CurrentUnits)
	}
	if trade.TakeProfitOrder == nil {
		t.Error("Expected TakeProfitOrder to be set")
	}
	if trade.StopLossOrder == nil {
		t.Error("Expected StopLossOrder to be set")
	}
	if trade.TrailingStopLossOrder == nil {
		t.Error("Expected TrailingStopLossOrder to be set")
	}
}

func TestGetOpenTrades(t *testing.T) {
	defer logTestResult(t, "TestGetOpenTrades")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/openTrades" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := ReceivedTrades{
			LastTransactionID: "1234",
			Trades: []Trade{
				{
					ID:                    "1",
					Instrument:            "EUR_USD",
					Price:                 "1.1000",
					OpenTime:              time.Now(),
					State:                 "OPEN",
					InitialUnits:          "100",
					InitialMarginRequired: "10.00",
					CurrentUnits:          "100",
					RealizedPL:            "0.00",
					UnrealizedPL:          "10.00",
					MarginUsed:            "10.00",
					AverageClosePrice:     "0.00",
					ClosingTransactionIDs: []string{},
					Financing:             "0.00",
					CloseTime:             time.Time{},
					ClientExtensions:      &OrderExtensions{Comment: "Test trade", Tag: "test"},
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

	trades, err := c.GetOpenTrades()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if trades.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", trades.LastTransactionID)
	}

	if len(trades.Trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(trades.Trades))
	}

	trade := trades.Trades[0]
	if trade.State != "OPEN" {
		t.Errorf("Expected State to be OPEN, got %s", trade.State)
	}
}

func TestGetTrade(t *testing.T) {
	defer logTestResult(t, "TestGetTrade")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/trades/1" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := ReceivedTrade{
			LastTransactionID: "1234",
			Trade: Trade{
				ID:                    "1",
				Instrument:            "EUR_USD",
				Price:                 "1.1000",
				OpenTime:              time.Now(),
				State:                 "OPEN",
				InitialUnits:          "100",
				InitialMarginRequired: "10.00",
				CurrentUnits:          "100",
				RealizedPL:            "0.00",
				UnrealizedPL:          "10.00",
				MarginUsed:            "10.00",
				AverageClosePrice:     "0.00",
				ClosingTransactionIDs: []string{},
				Financing:             "0.00",
				CloseTime:             time.Time{},
				ClientExtensions:      &OrderExtensions{Comment: "Test trade", Tag: "test"},
				TakeProfitOrder: &TakeProfitOrder{
					ID:               "2",
					CreateTime:       time.Now(),
					Type:             "TAKE_PROFIT",
					TradeID:          "1",
					Price:            "1.1100",
					TimeInForce:      "GTC",
					TriggerCondition: "DEFAULT",
					State:            "PENDING",
				},
				StopLossOrder: &StopLossOrder{
					ID:               "3",
					CreateTime:       time.Now(),
					Type:             "STOP_LOSS",
					TradeID:          "1",
					Price:            "1.0900",
					TimeInForce:      "GTC",
					TriggerCondition: "DEFAULT",
					State:            "PENDING",
					Guaranteed:       false,
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

	trade, err := c.GetTrade("1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if trade.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", trade.LastTransactionID)
	}

	if trade.Trade.ID != "1" {
		t.Errorf("Expected Trade ID to be 1, got %s", trade.Trade.ID)
	}

	if trade.Trade.TakeProfitOrder == nil {
		t.Error("Expected TakeProfitOrder to be set")
	}

	if trade.Trade.StopLossOrder == nil {
		t.Error("Expected StopLossOrder to be set")
	}
}

func TestReduceTradeSize(t *testing.T) {
	defer logTestResult(t, "TestReduceTradeSize")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/trades/1/close" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		var payload CloseTradePayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if payload.Units != "50" {
			t.Errorf("Expected Units to be 50, got %s", payload.Units)
			http.Error(w, "Invalid units", http.StatusBadRequest)
			return
		}

		response := ModifiedTrade{
			OrderCreateTransaction: struct {
				Type         string `json:"type"`
				Instrument   string `json:"instrument"`
				Units        string `json:"units"`
				TimeInForce  string `json:"timeInForce"`
				PositionFill string `json:"positionFill"`
				Reason       string `json:"reason"`
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
				Units:      "-50",
				TradeClose: struct {
					Units   string `json:"units"`
					TradeID string `json:"tradeID"`
				}{
					Units:   "50",
					TradeID: "1",
				},
				ID:        "1000",
				AccountID: "test-account",
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
				Units:      "-50",
				Price:      "1.1000",
				TradeReduced: struct {
					TradeID    string `json:"tradeID"`
					Units      string `json:"units"`
					RealizedPL string `json:"realizedPL"`
					Financing  string `json:"financing"`
				}{
					TradeID:    "1",
					Units:      "50",
					RealizedPL: "5.00",
					Financing:  "0.00",
				},
				ID:        "1001",
				AccountID: "test-account",
			},
			LastTransactionID: "1001",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "test-account",
		client:    *server.Client(),
	}

	modifiedTrade, err := c.ReduceTradeSize("1", CloseTradePayload{Units: "50"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if modifiedTrade.LastTransactionID != "1001" {
		t.Errorf("Expected LastTransactionID to be 1001, got %s", modifiedTrade.LastTransactionID)
	}

	if modifiedTrade.OrderCreateTransaction.Type != "MARKET_ORDER" {
		t.Errorf("Expected OrderCreateTransaction.Type to be MARKET_ORDER, got %s", modifiedTrade.OrderCreateTransaction.Type)
	}

	if modifiedTrade.OrderFillTransaction.TradeReduced.Units != "50" {
		t.Errorf("Expected OrderFillTransaction.TradeReduced.Units to be 50, got %s", modifiedTrade.OrderFillTransaction.TradeReduced.Units)
	}
}
