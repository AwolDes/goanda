package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAccounts(t *testing.T) {
	defer logTestResult(t, "Accounts")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// The API is expected to return an array of AccountProperties directly
		response := []AccountProperties{
			{
				ID:           "001-001-1234567-001",
				Mt4AccountID: 1234567,
				Tags:         []string{"demo", "test"},
			},
			{
				ID:           "001-001-1234568-001",
				Mt4AccountID: 1234568,
				Tags:         []string{"live"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	accounts, err := c.Accounts()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(accounts) != 2 {
		t.Fatalf("Expected 2 accounts, got %d", len(accounts))
	}

	if accounts[0].ID != "001-001-1234567-001" {
		t.Errorf("Expected first account ID to be 001-001-1234567-001, got %s", accounts[0].ID)
	}

	if accounts[1].Mt4AccountID != 1234568 {
		t.Errorf("Expected second account Mt4AccountID to be 1234568, got %d", accounts[1].Mt4AccountID)
	}
}

func TestGetAccount(t *testing.T) {
	defer logTestResult(t, "GetAccount")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/001-001-1234567-001" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := AccountInfo{
			Account: struct {
				NAV                         string        `json:"NAV"`
				Alias                       string        `json:"alias"`
				Balance                     string        `json:"balance"`
				CreatedByUserID             int           `json:"createdByUserID"`
				CreatedTime                 time.Time     `json:"createdTime"`
				Currency                    string        `json:"currency"`
				HedgingEnabled              bool          `json:"hedgingEnabled"`
				ID                          string        `json:"id"`
				LastTransactionID           string        `json:"lastTransactionID"`
				MarginAvailable             string        `json:"marginAvailable"`
				MarginCloseoutMarginUsed    string        `json:"marginCloseoutMarginUsed"`
				MarginCloseoutNAV           string        `json:"marginCloseoutNAV"`
				MarginCloseoutPercent       string        `json:"marginCloseoutPercent"`
				MarginCloseoutPositionValue string        `json:"marginCloseoutPositionValue"`
				MarginCloseoutUnrealizedPL  string        `json:"marginCloseoutUnrealizedPL"`
				MarginRate                  string        `json:"marginRate"`
				MarginUsed                  string        `json:"marginUsed"`
				OpenPositionCount           int           `json:"openPositionCount"`
				OpenTradeCount              int           `json:"openTradeCount"`
				Orders                      []interface{} `json:"orders"`
				PendingOrderCount           int           `json:"pendingOrderCount"`
				Pl                          string        `json:"pl"`
				PositionValue               string        `json:"positionValue"`
				Positions                   []struct {
					Instrument string `json:"instrument"`
					Long       struct {
						Pl           string `json:"pl"`
						ResettablePL string `json:"resettablePL"`
						Units        string `json:"units"`
						UnrealizedPL string `json:"unrealizedPL"`
					} `json:"long"`
					Pl           string `json:"pl"`
					ResettablePL string `json:"resettablePL"`
					Short        struct {
						Pl           string `json:"pl"`
						ResettablePL string `json:"resettablePL"`
						Units        string `json:"units"`
						UnrealizedPL string `json:"unrealizedPL"`
					} `json:"short"`
					UnrealizedPL string `json:"unrealizedPL"`
				} `json:"positions"`
				ResettablePL    string        `json:"resettablePL"`
				Trades          []interface{} `json:"trades"`
				UnrealizedPL    string        `json:"unrealizedPL"`
				WithdrawalLimit string        `json:"withdrawalLimit"`
			}{
				ID:                "001-001-1234567-001",
				NAV:               "43650.78",
				Balance:           "43650.78",
				Currency:          "USD",
				OpenPositionCount: 0,
				OpenTradeCount:    0,
				PendingOrderCount: 0,
			},
			LastTransactionID: "1234",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	account, err := c.GetAccount("001-001-1234567-001")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if account.Account.ID != "001-001-1234567-001" {
		t.Errorf("Expected account ID to be 001-001-1234567-001, got %s", account.Account.ID)
	}

	if account.Account.Balance != "43650.78" {
		t.Errorf("Expected account balance to be 43650.78, got %s", account.Account.Balance)
	}

	if account.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", account.LastTransactionID)
	}
}

func TestGetOrderDetails(t *testing.T) {
	defer logTestResult(t, "GetOrderDetails")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/001-001-1234567-001/orderEntryData" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		instrument := r.URL.Query().Get("instrument")
		if instrument != "EUR_USD" {
			t.Errorf("Expected instrument to be EUR_USD, got %s", instrument)
		}

		units := r.URL.Query().Get("units")
		if units != "100" {
			t.Errorf("Expected units to be 100, got %s", units)
		}

		response := OrderDetails{
			GainPerPipPerMillionUnits: 10.0,
			LossPerPipPerMillionUnits: 10.0,
			UnitsAvailable: struct {
				Default struct {
					Long  float64 `json:"long,string"`
					Short float64 `json:"short,string"`
				} `json:"default"`
				OpenOnly struct {
					Long  float64 `json:"long,string"`
					Short float64 `json:"short,string"`
				} `json:"openOnly"`
				ReduceFirst struct {
					Long  float64 `json:"long,string"`
					Short float64 `json:"short,string"`
				} `json:"reduceFirst"`
				ReduceOnly struct {
					Long  float64 `json:"long,string"`
					Short float64 `json:"short,string"`
				} `json:"reduceOnly"`
			}{
				Default: struct {
					Long  float64 `json:"long,string"`
					Short float64 `json:"short,string"`
				}{
					Long:  1000000,
					Short: 1000000,
				},
			},
			LastTransactionID: "1234",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "001-001-1234567-001",
		client:    *server.Client(),
	}

	details, err := c.GetOrderDetails("EUR_USD", "100")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if details.GainPerPipPerMillionUnits != 10.0 {
		t.Errorf("Expected GainPerPipPerMillionUnits to be 10.0, got %f", details.GainPerPipPerMillionUnits)
	}

	if details.UnitsAvailable.Default.Long != 1000000 {
		t.Errorf("Expected UnitsAvailable.Default.Long to be 1000000, got %f", details.UnitsAvailable.Default.Long)
	}

	if details.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", details.LastTransactionID)
	}
}

func TestGetAccountSummary(t *testing.T) {
	defer logTestResult(t, "GetAccountSummary")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/001-001-1234567-001/summary" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := AccountSummary{
			Account: struct {
				NAV                         string    `json:"NAV"`
				Alias                       string    `json:"alias"`
				Balance                     float64   `json:"balance,string"`
				CreatedByUserID             int       `json:"createdByUserID"`
				CreatedTime                 time.Time `json:"createdTime"`
				Currency                    string    `json:"currency"`
				HedgingEnabled              bool      `json:"hedgingEnabled"`
				ID                          string    `json:"id"`
				LastTransactionID           string    `json:"lastTransactionID"`
				MarginAvailable             float64   `json:"marginAvailable,string"`
				MarginCloseoutMarginUsed    string    `json:"marginCloseoutMarginUsed"`
				MarginCloseoutNAV           string    `json:"marginCloseoutNAV"`
				MarginCloseoutPercent       string    `json:"marginCloseoutPercent"`
				MarginCloseoutPositionValue string    `json:"marginCloseoutPositionValue"`
				MarginCloseoutUnrealizedPL  string    `json:"marginCloseoutUnrealizedPL"`
				MarginRate                  string    `json:"marginRate"`
				MarginUsed                  string    `json:"marginUsed"`
				OpenPositionCount           int       `json:"openPositionCount"`
				OpenTradeCount              int       `json:"openTradeCount"`
				PendingOrderCount           int       `json:"pendingOrderCount"`
				Pl                          string    `json:"pl"`
				PositionValue               string    `json:"positionValue"`
				ResettablePL                string    `json:"resettablePL"`
				UnrealizedPL                string    `json:"unrealizedPL"`
				WithdrawalLimit             string    `json:"withdrawalLimit"`
			}{
				ID:              "001-001-1234567-001",
				NAV:             "43650.78",
				Balance:         43650.78,
				Currency:        "USD",
				MarginAvailable: 43650.78,
			},
			LastTransactionID: "1234",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "001-001-1234567-001",
		client:    *server.Client(),
	}

	summary, err := c.GetAccountSummary()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if summary.Account.ID != "001-001-1234567-001" {
		t.Errorf("Expected account ID to be 001-001-1234567-001, got %s", summary.Account.ID)
	}

	if summary.Account.Balance != 43650.78 {
		t.Errorf("Expected account balance to be 43650.78, got %f", summary.Account.Balance)
	}

	if summary.LastTransactionID != "1234" {
		t.Errorf("Expected LastTransactionID to be 1234, got %s", summary.LastTransactionID)
	}
}

func TestGetAccountInstruments(t *testing.T) {
	defer logTestResult(t, "GetAccountInstruments")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/001-001-1234567-001/instruments" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := struct {
			Instruments Instruments `json:"instruments"`
		}{
			Instruments: Instruments{
				{
					DisplayName:                 "EUR/USD",
					DisplayPrecision:            5,
					MarginRate:                  "0.02",
					MaximumOrderUnits:           "100000000",
					MaximumPositionSize:         "0",
					MaximumTrailingStopDistance: "1",
					MinimumTradeSize:            "1",
					MinimumTrailingStopDistance: "0.0005",
					Name:                        "EUR_USD",
					PipLocation:                 -4,
					TradeUnitsPrecision:         0,
					Type:                        "CURRENCY",
				},
			},
		}
		responseJSON, _ := json.Marshal(response)
		t.Logf("Server sending response: %s", string(responseJSON))
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	instruments, err := c.GetAccountInstruments("001-001-1234567-001")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Logf("Received instruments: %+v", instruments)

	if len(instruments) != 1 {
		t.Fatalf("Expected 1 instrument, got %d", len(instruments))
	}

	if instruments[0].Name != "EUR_USD" {
		t.Errorf("Expected instrument name to be EUR_USD, got %s", instruments[0].Name)
	}

	if instruments[0].DisplayName != "EUR/USD" {
		t.Errorf("Expected instrument display name to be EUR/USD, got %s", instruments[0].DisplayName)
	}
}

func TestGetAccountChanges(t *testing.T) {
	defer logTestResult(t, "GetAccountChanges")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/001-001-1234567-001/changes" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		sinceTransactionID := r.URL.Query().Get("sinceTransactionID")
		if sinceTransactionID != "1234" {
			t.Errorf("Expected sinceTransactionID to be 1234, got %s", sinceTransactionID)
		}

		response := AccountChanges{
			Changes: struct {
				OrdersCancelled []interface{} `json:"ordersCancelled"`
				OrdersCreated   []interface{} `json:"ordersCreated"`
				OrdersFilled    []struct {
					CreateTime           time.Time `json:"createTime"`
					FilledTime           time.Time `json:"filledTime"`
					FillingTransactionID string    `json:"fillingTransactionID"`
					ID                   string    `json:"id"`
					Instrument           string    `json:"instrument"`
					PositionFill         string    `json:"positionFill"`
					State                string    `json:"state"`
					TimeInForce          string    `json:"timeInForce"`
					TradeOpenedID        string    `json:"tradeOpenedID"`
					Type                 string    `json:"type"`
					Units                string    `json:"units"`
				} `json:"ordersFilled"`
				OrdersTriggered []interface{} `json:"ordersTriggered"`
				Positions       []struct {
					Instrument string `json:"instrument"`
					Long       struct {
						Pl           string `json:"pl"`
						ResettablePL string `json:"resettablePL"`
						Units        string `json:"units"`
					} `json:"long"`
					Pl           string `json:"pl"`
					ResettablePL string `json:"resettablePL"`
					Short        struct {
						AveragePrice string   `json:"averagePrice"`
						Pl           string   `json:"pl"`
						ResettablePL string   `json:"resettablePL"`
						TradeIDs     []string `json:"tradeIDs"`
						Units        string   `json:"units"`
					} `json:"short"`
				} `json:"positions"`
				TradesClosed []interface{} `json:"tradesClosed"`
				TradesOpened []struct {
					CurrentUnits string    `json:"currentUnits"`
					Financing    string    `json:"financing"`
					ID           string    `json:"id"`
					InitialUnits string    `json:"initialUnits"`
					Instrument   string    `json:"instrument"`
					OpenTime     time.Time `json:"openTime"`
					Price        string    `json:"price"`
					RealizedPL   string    `json:"realizedPL"`
					State        string    `json:"state"`
				} `json:"tradesOpened"`
				TradesReduced []interface{} `json:"tradesReduced"`
				Transactions  []struct {
					AccountID      string    `json:"accountID"`
					BatchID        string    `json:"batchID"`
					ID             string    `json:"id"`
					Instrument     string    `json:"instrument"`
					PositionFill   string    `json:"positionFill,omitempty"`
					Reason         string    `json:"reason"`
					Time           time.Time `json:"time"`
					TimeInForce    string    `json:"timeInForce,omitempty"`
					Type           string    `json:"type"`
					Units          string    `json:"units"`
					UserID         int       `json:"userID"`
					AccountBalance string    `json:"accountBalance,omitempty"`
					Financing      string    `json:"financing,omitempty"`
					OrderID        string    `json:"orderID,omitempty"`
					Pl             string    `json:"pl,omitempty"`
					Price          string    `json:"price,omitempty"`
					TradeOpened    struct {
						TradeID string `json:"tradeID"`
						Units   string `json:"units"`
					} `json:"tradeOpened,omitempty"`
				} `json:"transactions"`
			}{}, // Add this comma
			LastTransactionID: "1235",
			State: struct {
				NAV                        string        `json:"NAV"`
				MarginAvailable            string        `json:"marginAvailable"`
				MarginCloseoutMarginUsed   string        `json:"marginCloseoutMarginUsed"`
				MarginCloseoutNAV          string        `json:"marginCloseoutNAV"`
				MarginCloseoutPercent      string        `json:"marginCloseoutPercent"`
				MarginCloseoutUnrealizedPL string        `json:"marginCloseoutUnrealizedPL"`
				MarginUsed                 string        `json:"marginUsed"`
				Orders                     []interface{} `json:"orders"`
				PositionValue              string        `json:"positionValue"`
				Positions                  []struct {
					Instrument        string `json:"instrument"`
					LongUnrealizedPL  string `json:"longUnrealizedPL"`
					NetUnrealizedPL   string `json:"netUnrealizedPL"`
					ShortUnrealizedPL string `json:"shortUnrealizedPL"`
				} `json:"positions"`
				Trades []struct {
					ID           string `json:"id"`
					UnrealizedPL string `json:"unrealizedPL"`
				} `json:"trades"`
				UnrealizedPL    string `json:"unrealizedPL"`
				WithdrawalLimit string `json:"withdrawalLimit"`
			}{
				NAV:             "10000.00",
				MarginAvailable: "9000.00",
			},
		}

		// Populate the response with some sample data
		response.LastTransactionID = "1235"
		response.State.NAV = "10000.00"
		response.State.MarginAvailable = "9000.00"

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	changes, err := c.GetAccountChanges("001-001-1234567-001", "1234")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if changes.LastTransactionID != "1235" {
		t.Errorf("Expected LastTransactionID to be 1235, got %s", changes.LastTransactionID)
	}

	if changes.State.NAV != "10000.00" {
		t.Errorf("Expected NAV to be 10000.00, got %s", changes.State.NAV)
	}

	if changes.State.MarginAvailable != "9000.00" {
		t.Errorf("Expected MarginAvailable to be 9000.00, got %s", changes.State.MarginAvailable)
	}
}
