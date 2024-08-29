package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPricingForInstruments(t *testing.T) {
	defer logTestResult(t, "GetPricingForInstruments")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/test-account/pricing" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		instruments := r.URL.Query().Get("instruments")
		if instruments != "EUR_USD,USD_JPY" {
			t.Errorf("Unexpected instruments: %s", instruments)
			http.Error(w, "Invalid instruments", http.StatusBadRequest)
			return
		}

		response := Pricings{
			Prices: []struct {
				Asks []struct {
					Liquidity int    `json:"liquidity"`
					Price     string `json:"price"`
				} `json:"asks"`
				Bids []struct {
					Liquidity int    `json:"liquidity"`
					Price     string `json:"price"`
				} `json:"bids"`
				CloseoutAsk                string `json:"closeoutAsk"`
				CloseoutBid                string `json:"closeoutBid"`
				Instrument                 string `json:"instrument"`
				QuoteHomeConversionFactors struct {
					NegativeUnits string `json:"negativeUnits"`
					PositiveUnits string `json:"positiveUnits"`
				} `json:"quoteHomeConversionFactors"`
				Status         string    `json:"status"`
				Time           time.Time `json:"time"`
				UnitsAvailable struct {
					Default struct {
						Long  string `json:"long"`
						Short string `json:"short"`
					} `json:"default"`
					OpenOnly struct {
						Long  string `json:"long"`
						Short string `json:"short"`
					} `json:"openOnly"`
					ReduceFirst struct {
						Long  string `json:"long"`
						Short string `json:"short"`
					} `json:"reduceFirst"`
					ReduceOnly struct {
						Long  string `json:"long"`
						Short string `json:"short"`
					} `json:"reduceOnly"`
				} `json:"unitsAvailable"`
			}{
				{
					Asks: []struct {
						Liquidity int    `json:"liquidity"`
						Price     string `json:"price"`
					}{
						{Liquidity: 10000000, Price: "1.10050"},
					},
					Bids: []struct {
						Liquidity int    `json:"liquidity"`
						Price     string `json:"price"`
					}{
						{Liquidity: 10000000, Price: "1.10040"},
					},
					CloseoutAsk: "1.10060",
					CloseoutBid: "1.10030",
					Instrument:  "EUR_USD",
					Status:      "tradeable",
					Time:        time.Now(),
				},
				{
					Asks: []struct {
						Liquidity int    `json:"liquidity"`
						Price     string `json:"price"`
					}{
						{Liquidity: 10000000, Price: "109.500"},
					},
					Bids: []struct {
						Liquidity int    `json:"liquidity"`
						Price     string `json:"price"`
					}{
						{Liquidity: 10000000, Price: "109.490"},
					},
					CloseoutAsk: "109.510",
					CloseoutBid: "109.480",
					Instrument:  "USD_JPY",
					Status:      "tradeable",
					Time:        time.Now(),
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

	pricings, err := c.GetPricingForInstruments([]string{"EUR_USD", "USD_JPY"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pricings.Prices) != 2 {
		t.Fatalf("Expected 2 prices, got %d", len(pricings.Prices))
	}

	eurUsd := pricings.Prices[0]
	if eurUsd.Instrument != "EUR_USD" {
		t.Errorf("Expected first instrument to be EUR_USD, got %s", eurUsd.Instrument)
	}
	if eurUsd.Status != "tradeable" {
		t.Errorf("Expected EUR_USD status to be tradeable, got %s", eurUsd.Status)
	}
	if len(eurUsd.Asks) == 0 || eurUsd.Asks[0].Price != "1.10050" {
		t.Errorf("Unexpected EUR_USD ask price: %v", eurUsd.Asks)
	}
	if len(eurUsd.Bids) == 0 || eurUsd.Bids[0].Price != "1.10040" {
		t.Errorf("Unexpected EUR_USD bid price: %v", eurUsd.Bids)
	}

	usdJpy := pricings.Prices[1]
	if usdJpy.Instrument != "USD_JPY" {
		t.Errorf("Expected second instrument to be USD_JPY, got %s", usdJpy.Instrument)
	}
	if usdJpy.Status != "tradeable" {
		t.Errorf("Expected USD_JPY status to be tradeable, got %s", usdJpy.Status)
	}
	if len(usdJpy.Asks) == 0 || usdJpy.Asks[0].Price != "109.500" {
		t.Errorf("Unexpected USD_JPY ask price: %v", usdJpy.Asks)
	}
	if len(usdJpy.Bids) == 0 || usdJpy.Bids[0].Price != "109.490" {
		t.Errorf("Unexpected USD_JPY bid price: %v", usdJpy.Bids)
	}
}
