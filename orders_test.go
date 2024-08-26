package goanda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/accounts/test-account/orders", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var payload OrderPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		units := strconv.Itoa(payload.Order.Units)

		response := OrderResponse{
			LastTransactionID: "1000",
			OrderCreateTransaction: struct {
				AccountID                string           `json:"accountID"`
				BatchID                  string           `json:"batchID"`
				ID                       string           `json:"id"`
				Instrument               string           `json:"instrument"`
				PositionFill             string           `json:"positionFill"`
				Reason                   string           `json:"reason"`
				Time                     time.Time        `json:"time"`
				TimeInForce              string           `json:"timeInForce"`
				Type                     string           `json:"type"`
				Units                    string           `json:"units"`
				UserID                   int              `json:"userID"`
				Price                    string           `json:"price,omitempty"`
				PriceBound               string           `json:"priceBound,omitempty"`
				Extensions               *OrderExtensions `json:"clientExtensions,omitempty"`
				TakeProfitOnFill         *OnFill          `json:"takeProfitOnFill,omitempty"`
				StopLossOnFill           *OnFill          `json:"stopLossOnFill,omitempty"`
				GuaranteedStopLossOnFill *OnFill          `json:"guaranteedStopLossOnFill,omitempty"`
				TrailingStopLossOnFill   *OnFill          `json:"trailingStopLossOnFill,omitempty"`
				TradeClientExtensions    *OrderExtensions `json:"tradeClientExtensions,omitempty"`
				TriggerCondition         string           `json:"triggerCondition,omitempty"`
				GTDTime                  time.Time        `json:"gtdTime,omitempty"`
				Distance                 string           `json:"distance,omitempty"`
			}{
				ID:           "1000",
				Type:         "MARKET_ORDER",
				Instrument:   payload.Order.Instrument,
				Units:        units,
				TimeInForce:  payload.Order.TimeInForce,
				PositionFill: payload.Order.PositionFill,
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

	order := OrderPayload{
		Order: OrderBody{
			Type:         "MARKET",
			Instrument:   "EUR_USD",
			Units:        100,
			TimeInForce:  "FOK",
			PositionFill: "DEFAULT",
			TakeProfitOnFill: &OnFill{
				Price: "1.25000",
			},
			StopLossOnFill: &OnFill{
				Price: "1.20000",
			},
			TrailingStopLossOnFill: &OnFill{
				Distance: "0.01000",
			},
			ClientExtensions: &OrderExtensions{
				Comment: "Test order",
				Tag:     "test",
			},
		},
	}

	response, err := c.CreateOrder(order)
	assert.NoError(t, err)

	assert.Equal(t, "1000", response.LastTransactionID)
	assert.Equal(t, "MARKET_ORDER", response.OrderCreateTransaction.Type)
	assert.Equal(t, "EUR_USD", response.OrderCreateTransaction.Instrument)
	assert.Equal(t, "100", response.OrderCreateTransaction.Units)
	assert.Equal(t, "FOK", response.OrderCreateTransaction.TimeInForce)
	assert.Equal(t, "DEFAULT", response.OrderCreateTransaction.PositionFill)
}

func TestGetOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/accounts/test-account/orders", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		instrument := r.URL.Query().Get("instrument")
		assert.Equal(t, "EUR_USD", instrument)

		response := RetrievedOrders{
			LastTransactionID: "1000",
			Orders: []OrderInfo{
				{
					ID:               "1",
					CreateTime:       time.Now(),
					Type:             "LIMIT",
					Instrument:       "EUR_USD",
					Units:            "100",
					Price:            "1.1000",
					State:            "PENDING",
					TimeInForce:      "GTC",
					PositionFill:     "DEFAULT",
					TriggerCondition: "DEFAULT",
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

	orders, err := c.GetOrders("EUR_USD")
	assert.NoError(t, err)

	assert.Equal(t, "1000", orders.LastTransactionID)
	assert.Len(t, orders.Orders, 1)

	order := orders.Orders[0]
	assert.Equal(t, "1", order.ID)
	assert.Equal(t, "LIMIT", order.Type)
	assert.Equal(t, "EUR_USD", order.Instrument)
	assert.Equal(t, "100", order.Units)
	assert.Equal(t, "1.1000", order.Price)
	assert.Equal(t, "PENDING", order.State)
	assert.Equal(t, "GTC", order.TimeInForce)
	assert.Equal(t, "DEFAULT", order.PositionFill)
	assert.Equal(t, "DEFAULT", order.TriggerCondition)
}

func TestGetOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/accounts/test-account/orders/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		response := RetrievedOrder{
			Order: OrderInfo{
				ID:               "1",
				CreateTime:       time.Now(),
				Type:             "LIMIT",
				Instrument:       "EUR_USD",
				Units:            "100",
				Price:            "1.1000",
				State:            "PENDING",
				TimeInForce:      "GTC",
				PositionFill:     "DEFAULT",
				TriggerCondition: "DEFAULT",
				TakeProfitOnFill: &OnFill{
					Price: "1.2000",
				},
				StopLossOnFill: &OnFill{
					Price: "1.0000",
				},
				TrailingStopLossOnFill: &OnFill{
					Distance: "0.0100",
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

	order, err := c.GetOrder("1")
	assert.NoError(t, err)

	assert.Equal(t, "1", order.Order.ID)
	assert.Equal(t, "LIMIT", order.Order.Type)
	assert.Equal(t, "EUR_USD", order.Order.Instrument)
	assert.Equal(t, "100", order.Order.Units)
	assert.Equal(t, "1.1000", order.Order.Price)
	assert.Equal(t, "PENDING", order.Order.State)
	assert.Equal(t, "GTC", order.Order.TimeInForce)
	assert.Equal(t, "DEFAULT", order.Order.PositionFill)
	assert.Equal(t, "DEFAULT", order.Order.TriggerCondition)
	assert.Equal(t, "1.2000", order.Order.TakeProfitOnFill.Price)
	assert.Equal(t, "1.0000", order.Order.StopLossOnFill.Price)
	assert.Equal(t, "0.0100", order.Order.TrailingStopLossOnFill.Distance)
}

func TestUpdateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/accounts/test-account/orders/1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		var payload OrderPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		units := strconv.Itoa(payload.Order.Units)

		response := RetrievedOrder{
			Order: OrderInfo{
				ID:                     "1",
				CreateTime:             time.Now(),
				Type:                   payload.Order.Type,
				Instrument:             payload.Order.Instrument,
				Units:                  units,
				Price:                  payload.Order.Price,
				State:                  "PENDING",
				TimeInForce:            payload.Order.TimeInForce,
				PositionFill:           payload.Order.PositionFill,
				TriggerCondition:       "DEFAULT",
				TakeProfitOnFill:       payload.Order.TakeProfitOnFill,
				StopLossOnFill:         payload.Order.StopLossOnFill,
				TrailingStopLossOnFill: payload.Order.TrailingStopLossOnFill,
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

	updatedOrder := OrderPayload{
		Order: OrderBody{
			Type:         "LIMIT",
			Instrument:   "EUR_USD",
			Units:        200,
			Price:        "1.1100",
			TimeInForce:  "GTC",
			PositionFill: "DEFAULT",
			TakeProfitOnFill: &OnFill{
				Price: "1.2100",
			},
			StopLossOnFill: &OnFill{
				Price: "1.0100",
			},
			TrailingStopLossOnFill: &OnFill{
				Distance: "0.0200",
			},
		},
	}

	order, err := c.UpdateOrder("1", updatedOrder)
	assert.NoError(t, err)

	assert.Equal(t, "1", order.Order.ID)
	assert.Equal(t, "LIMIT", order.Order.Type)
	assert.Equal(t, "EUR_USD", order.Order.Instrument)
	assert.Equal(t, "200", order.Order.Units)
	assert.Equal(t, "1.1100", order.Order.Price)
	assert.Equal(t, "PENDING", order.Order.State)
	assert.Equal(t, "GTC", order.Order.TimeInForce)
	assert.Equal(t, "DEFAULT", order.Order.PositionFill)
	assert.Equal(t, "DEFAULT", order.Order.TriggerCondition)
	assert.Equal(t, "1.2100", order.Order.TakeProfitOnFill.Price)
	assert.Equal(t, "1.0100", order.Order.StopLossOnFill.Price)
	assert.Equal(t, "0.0200", order.Order.TrailingStopLossOnFill.Distance)
}

func TestCancelOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/accounts/test-account/orders/1/cancel", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)

		response := CancelledOrder{
			OrderCancelTransaction: struct {
				ID                string    `json:"id"`
				Time              time.Time `json:"time"`
				UserID            int       `json:"userID"`
				AccountID         string    `json:"accountID"`
				BatchID           string    `json:"batchID"`
				RequestID         string    `json:"requestID"`
				Type              string    `json:"type"`
				OrderID           string    `json:"orderID"`
				ClientOrderID     string    `json:"clientOrderID"`
				Reason            string    `json:"reason"`
				ReplacedByOrderID string    `json:"replacedByOrderID"`
			}{
				ID:      "1000",
				Type:    "ORDER_CANCEL",
				OrderID: "1",
				Reason:  "CLIENT_REQUEST",
			},
			RelatedTransactionIDs: []string{"1000"},
			LastTransactionID:     "1000",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	c := &Connection{
		hostname:  server.URL,
		accountID: "test-account",
		client:    *server.Client(),
	}

	cancelledOrder, err := c.CancelOrder("1")
	assert.NoError(t, err)

	assert.Equal(t, "1000", cancelledOrder.OrderCancelTransaction.ID)
	assert.Equal(t, "ORDER_CANCEL", cancelledOrder.OrderCancelTransaction.Type)
	assert.Equal(t, "1", cancelledOrder.OrderCancelTransaction.OrderID)
	assert.Equal(t, "CLIENT_REQUEST", cancelledOrder.OrderCancelTransaction.Reason)
	assert.Equal(t, []string{"1000"}, cancelledOrder.RelatedTransactionIDs)
	assert.Equal(t, "1000", cancelledOrder.LastTransactionID)
}

func TestOrderIntegration(t *testing.T) {
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

	// Get current price for AUD_CAD
	instrument := "AUD_CAD"
	pricing, err := conn.GetPricingForInstruments([]string{instrument})
	if err != nil {
		t.Fatalf("Error getting pricing: %v", err)
	}
	if len(pricing.Prices) == 0 {
		t.Fatalf("No pricing information received for %s", instrument)
	}
	currentPrice, err := strconv.ParseFloat(pricing.Prices[0].CloseoutAsk, 64)
	if err != nil {
		t.Fatalf("Error parsing current price: %v", err)
	}

	// Calculate take profit and stop loss prices
	takeProfitPrice := fmt.Sprintf("%.5f", currentPrice*1.01) // 1% above current price
	stopLossPrice := fmt.Sprintf("%.5f", currentPrice*0.99)   // 1% below current price
	trailingStopDistance := "0.00500"                         // 50 pips

	// Create a new order
	orderBody := OrderBody{
		Instrument:   instrument,
		Units:        100,
		Type:         "MARKET",
		TimeInForce:  "FOK",
		PositionFill: "DEFAULT",
		ClientExtensions: &OrderExtensions{
			Comment: "Test order",
			Tag:     "integration-test",
		},
		TakeProfitOnFill: &OnFill{
			Price: takeProfitPrice,
		},
		StopLossOnFill: &OnFill{
			Price: stopLossPrice,
		},
		TrailingStopLossOnFill: &OnFill{
			Distance: trailingStopDistance,
		},
		TradeClientExtensions: &OrderExtensions{
			Comment: "Test trade",
			Tag:     "integration-test-trade",
		},
	}

	orderResponse, err := conn.CreateOrder(OrderPayload{Order: orderBody})
	if err != nil {
		t.Fatalf("Error creating order: %v", err)
	}

	// Verify the order creation response
	if orderResponse.OrderCreateTransaction.Type != "MARKET_ORDER" {
		t.Errorf("Expected order type MARKET_ORDER, got %s", orderResponse.OrderCreateTransaction.Type)
	}
	if orderResponse.OrderCreateTransaction.Instrument != instrument {
		t.Errorf("Expected instrument %s, got %s", instrument, orderResponse.OrderCreateTransaction.Instrument)
	}

	// Retrieve the created order
	time.Sleep(2 * time.Second) // Wait for order to be processed
	retrievedOrder, err := conn.GetOrder(orderResponse.OrderCreateTransaction.ID)
	if err != nil {
		t.Fatalf("Error retrieving order: %v", err)
	}

	// Verify order details
	order := retrievedOrder.Order
	if order.Instrument != instrument {
		t.Errorf("Expected instrument %s, got %s", instrument, order.Instrument)
	}
	if order.Units != "100" {
		t.Errorf("Expected units 100, got %s", order.Units)
	}
	if order.Type != "MARKET" {
		t.Errorf("Expected type MARKET, got %s", order.Type)
	}
	if order.State != "FILLED" {
		t.Errorf("Expected state FILLED, got %s", order.State)
	}
	if order.ClientExtensions.Comment != "Test order" {
		t.Errorf("Expected client extension comment 'Test order', got %s", order.ClientExtensions.Comment)
	}
	if order.TradeOpenedID == "" {
		t.Error("Expected TradeOpenedID to be set")
	}

	// Verify take profit, stop loss, and trailing stop loss orders
	trades, err := conn.GetTradesForInstrument(instrument)
	if err != nil {
		t.Fatalf("Error retrieving trades: %v", err)
	}

	if len(trades.Trades) == 0 {
		t.Fatal("Expected at least one trade")
	}

	trade := trades.Trades[0]

	// Verify take profit, stop loss, and trailing stop loss
	if trade.TakeProfitOrder == nil || trade.TakeProfitOrder.Price != takeProfitPrice {
		t.Errorf("Expected take profit price %s, got %v", takeProfitPrice, trade.TakeProfitOrder)
	}
	if trade.StopLossOrder == nil || trade.StopLossOrder.Price != stopLossPrice {
		t.Errorf("Expected stop loss price %s, got %v", stopLossPrice, trade.StopLossOrder)
	}
	if trade.TrailingStopLossOrder == nil || trade.TrailingStopLossOrder.Distance != trailingStopDistance {
		t.Errorf("Expected trailing stop loss distance %s, got %v", trailingStopDistance, trade.TrailingStopLossOrder)
	}

	// Clean up: close the position
	time.Sleep(5 * time.Second)
	_, err = conn.ReduceTradeSize(trade.ID, CloseTradePayload{Units: "ALL"})
	if err != nil {
		t.Fatalf("Error closing position: %v", err)
	}
}
