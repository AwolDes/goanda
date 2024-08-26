package goanda

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGranularityFromDuration(t *testing.T) {
	defer logTestResult(t, "GranularityFromDuration")

	tests := []struct {
		duration time.Duration
		expected Granularity
		hasError bool
	}{
		{time.Second * 5, GranularityFiveSeconds, false},
		{time.Minute, GranularityMinute, false},
		{time.Hour * 24, GranularityDay, false},
		{time.Hour * 24 * 7, GranularityWeek, false},
		{time.Second * 7, 0, true}, // Invalid granularity
	}

	for _, test := range tests {
		result, err := GranularityFromDuration(test.duration)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for duration %v, but got none", test.duration)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for duration %v: %v", test.duration, err)
			}
			if result != test.expected {
				t.Errorf("For duration %v, expected %v, but got %v", test.duration, test.expected, result)
			}
		}
	}
}

func TestGranularityDuration(t *testing.T) {
	defer logTestResult(t, "GranularityDuration")

	tests := []struct {
		granularity Granularity
		expected    time.Duration
	}{
		{GranularityFiveSeconds, time.Second * 5},
		{GranularityMinute, time.Minute},
		{GranularityDay, time.Hour * 24},
		{GranularityWeek, time.Hour * 24 * 7},
	}

	for _, test := range tests {
		result := test.granularity.Duration()
		if result != test.expected {
			t.Errorf("For granularity %v, expected %v, but got %v", test.granularity, test.expected, result)
		}
	}
}

func TestGranularityString(t *testing.T) {
	defer logTestResult(t, "GranularityString")

	tests := []struct {
		granularity Granularity
		expected    string
	}{
		{GranularityFiveSeconds, "S5"},
		{GranularityMinute, "M1"},
		{GranularityDay, "D"},
		{GranularityWeek, "W"},
	}

	for _, test := range tests {
		result := test.granularity.String()
		if result != test.expected {
			t.Errorf("For granularity %v, expected %v, but got %v", test.granularity, test.expected, result)
		}
	}
}

func TestGetCandles(t *testing.T) {
	defer logTestResult(t, "GetCandles")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instruments/EUR_USD/candles" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		count := r.URL.Query().Get("count")
		if count != "10" {
			t.Errorf("Expected count to be 10, got %s", count)
		}

		granularity := r.URL.Query().Get("granularity")
		if granularity != "M1" {
			t.Errorf("Expected granularity to be M1, got %s", granularity)
		}

		response := InstrumentHistory{
			Instrument:  "EUR_USD",
			Granularity: "M1",
			Candles: []Candles{
				{
					Complete: true,
					Volume:   100,
					Time:     time.Now(),
					Mid: Candle{
						Open:  1.1000,
						High:  1.1010,
						Low:   1.0990,
						Close: 1.1005,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	history, err := c.GetCandles("EUR_USD", 10, GranularityMinute)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if history.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", history.Instrument)
	}

	if history.Granularity != "M1" {
		t.Errorf("Expected Granularity to be M1, got %s", history.Granularity)
	}

	if len(history.Candles) != 1 {
		t.Fatalf("Expected 1 candle, got %d", len(history.Candles))
	}

	candle := history.Candles[0]
	if !candle.Complete {
		t.Errorf("Expected Complete to be true")
	}

	if candle.Volume != 100 {
		t.Errorf("Expected Volume to be 100, got %d", candle.Volume)
	}

	if candle.Mid.Open != 1.1000 {
		t.Errorf("Expected Open to be 1.1000, got %f", candle.Mid.Open)
	}
}

func TestGetTimeToCandles(t *testing.T) {
	defer logTestResult(t, "GetTimeToCandles")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instruments/EUR_USD/candles" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		count := r.URL.Query().Get("count")
		if count != "10" {
			t.Errorf("Expected count to be 10, got %s", count)
		}

		to := r.URL.Query().Get("to")
		if to == "" {
			t.Errorf("Expected 'to' parameter, but it's missing")
		}

		granularity := r.URL.Query().Get("granularity")
		if granularity != "M1" {
			t.Errorf("Expected granularity to be M1, got %s", granularity)
		}

		response := InstrumentHistory{
			Instrument:  "EUR_USD",
			Granularity: "M1",
			Candles: []Candles{
				{
					Complete: true,
					Volume:   100,
					Time:     time.Now(),
					Mid: Candle{
						Open:  1.1000,
						High:  1.1010,
						Low:   1.0990,
						Close: 1.1005,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	to := time.Now()
	history, err := c.GetTimeToCandles("EUR_USD", 10, GranularityMinute, to)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if history.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", history.Instrument)
	}

	if history.Granularity != "M1" {
		t.Errorf("Expected Granularity to be M1, got %s", history.Granularity)
	}

	if len(history.Candles) != 1 {
		t.Fatalf("Expected 1 candle, got %d", len(history.Candles))
	}
}

func TestGetTimeFromCandles(t *testing.T) {
	defer logTestResult(t, "GetTimeFromCandles")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instruments/EUR_USD/candles" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		count := r.URL.Query().Get("count")
		if count != "10" {
			t.Errorf("Expected count to be 10, got %s", count)
		}

		from := r.URL.Query().Get("from")
		if from == "" {
			t.Errorf("Expected 'from' parameter, but it's missing")
		}

		granularity := r.URL.Query().Get("granularity")
		if granularity != "M1" {
			t.Errorf("Expected granularity to be M1, got %s", granularity)
		}

		response := InstrumentHistory{
			Instrument:  "EUR_USD",
			Granularity: "M1",
			Candles: []Candles{
				{
					Complete: true,
					Volume:   100,
					Time:     time.Now(),
					Mid: Candle{
						Open:  1.1000,
						High:  1.1010,
						Low:   1.0990,
						Close: 1.1005,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	from := time.Now().Add(-time.Hour)
	history, err := c.GetTimeFromCandles("EUR_USD", 10, GranularityMinute, from)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if history.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", history.Instrument)
	}

	if history.Granularity != "M1" {
		t.Errorf("Expected Granularity to be M1, got %s", history.Granularity)
	}

	if len(history.Candles) != 1 {
		t.Fatalf("Expected 1 candle, got %d", len(history.Candles))
	}
}

func TestGetBidAskCandles(t *testing.T) {
	defer logTestResult(t, "GetBidAskCandles")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instruments/EUR_USD/candles" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		count := r.URL.Query().Get("count")
		if count != "10" {
			t.Errorf("Expected count to be 10, got %s", count)
		}

		granularity := r.URL.Query().Get("granularity")
		if granularity != "M1" {
			t.Errorf("Expected granularity to be M1, got %s", granularity)
		}

		price := r.URL.Query().Get("price")
		if price != "BA" {
			t.Errorf("Expected price to be BA, got %s", price)
		}

		response := BidAskCandles{
			Candles: []struct {
				Ask struct {
					C float64 `json:"c,string"`
					H float64 `json:"h,string"`
					L float64 `json:"l,string"`
					O float64 `json:"o,string"`
				} `json:"ask"`
				Bid struct {
					C float64 `json:"c,string"`
					H float64 `json:"h,string"`
					L float64 `json:"l,string"`
					O float64 `json:"o,string"`
				} `json:"bid"`
				Complete bool      `json:"complete"`
				Time     time.Time `json:"time"`
				Volume   int       `json:"volume"`
			}{
				{
					Ask: struct {
						C float64 `json:"c,string"`
						H float64 `json:"h,string"`
						L float64 `json:"l,string"`
						O float64 `json:"o,string"`
					}{
						O: 1.1000,
						H: 1.1010,
						L: 1.0990,
						C: 1.1005,
					},
					Bid: struct {
						C float64 `json:"c,string"`
						H float64 `json:"h,string"`
						L float64 `json:"l,string"`
						O float64 `json:"o,string"`
					}{
						O: 1.0998,
						H: 1.1008,
						L: 1.0988,
						C: 1.1003,
					},
					Complete: true,
					Time:     time.Now(),
					Volume:   100,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	candles, err := c.GetBidAskCandles("EUR_USD", "10", GranularityMinute)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(candles.Candles) != 1 {
		t.Fatalf("Expected 1 candle, got %d", len(candles.Candles))
	}

	candle := candles.Candles[0]
	if !candle.Complete {
		t.Errorf("Expected Complete to be true")
	}

	if candle.Volume != 100 {
		t.Errorf("Expected Volume to be 100, got %d", candle.Volume)
	}

	if candle.Ask.O != 1.1000 {
		t.Errorf("Expected Ask Open to be 1.1000, got %f", candle.Ask.O)
	}

	if candle.Bid.O != 1.0998 {
		t.Errorf("Expected Bid Open to be 1.0998, got %f", candle.Bid.O)
	}
}

func TestOrderBook(t *testing.T) {
	defer logTestResult(t, "OrderBook")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/instruments/EUR_USD/orderBook" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		response := BrokerBook{
			Instrument:  "EUR_USD",
			Time:        time.Now(),
			Price:       "1.1000",
			BucketWidth: "0.0001",
			Buckets: []Bucket{
				{
					Price:      "1.0999",
					LongCountPercent:  "40",
					ShortCountPercent: "60",
				},
				{
					Price:      "1.1000",
					LongCountPercent:  "50",
					ShortCountPercent: "50",
				},
				{
					Price:      "1.1001",
					LongCountPercent:  "60",
					ShortCountPercent: "40",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname: server.URL,
		client:   *server.Client(),
	}

	book, err := c.OrderBook("EUR_USD")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if book.Instrument != "EUR_USD" {
		t.Errorf("Expected Instrument to be EUR_USD, got %s", book.Instrument)
	}

	if book.Price != "1.1000" {
		t.Errorf("Expected Price to be 1.1000, got %s", book.Price)
	}

	if book.BucketWidth != "0.0001" {
		t.Errorf("Expected BucketWidth to be 0.0001, got %s", book.BucketWidth)
	}

	if len(book.Buckets) != 3 {
		t.Fatalf("Expected 3 buckets, got %d", len(book.Buckets))
	}

	if book.Buckets[0].Price != "1.0999" {
		t.Errorf("Expected first bucket Price to be 1.0999, got %s", book.Buckets[0].Price)
	}

	if book.Buckets[0].LongCountPercent != "40" {
		t.Errorf("Expected first bucket LongCountPercent to be 40, got %s", book.Buckets[0].LongCountPercent)
	}

	if book.Buckets[0].ShortCountPercent != "60" {
		t.Errorf("Expected first bucket ShortCountPercent to be 60, got %s", book.Buckets[0].ShortCountPercent)
	}
}