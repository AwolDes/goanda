package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/order-ep/

import (
	"time"
)

type OrderExtensions struct {
	Comment string `json:"comment,omitempty"`
	ID      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type OnFill struct {
	TimeInForce      string           `json:"timeInForce,omitempty"`
	Price            string           `json:"price,omitempty"` // must be a string for float precision
	Distance         string           `json:"distance,omitempty"`
	GtdTime          string           `json:"gtdTime,omitempty"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
}

type OrderBody struct {
	ID                       string           `json:"id,omitempty"`
	CreateTime               time.Time        `json:"createTime,omitempty"`
	State                    string           `json:"state,omitempty"`
	ClientExtensions         *OrderExtensions `json:"clientExtensions,omitempty"`
	Instrument               string           `json:"instrument"`
	Units                    int              `json:"units"`
	TimeInForce              string           `json:"timeInForce"`
	PriceBound               string           `json:"priceBound,omitempty"`
	Type                     string           `json:"type"`
	PositionFill             string           `json:"positionFill,omitempty"`
	Price                    string           `json:"price,omitempty"`
	TakeProfitOnFill         *OnFill          `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill           *OnFill          `json:"stopLossOnFill,omitempty"`
	GuaranteedStopLossOnFill *OnFill          `json:"guaranteedStopLossOnFill,omitempty"`
	TrailingStopLossOnFill   *OnFill          `json:"trailingStopLossOnFill,omitempty"`
	TradeClientExtensions    *OrderExtensions `json:"tradeClientExtensions,omitempty"`
	FillingTransactionID     string           `json:"fillingTransactionID,omitempty"`
	FilledTime               time.Time        `json:"filledTime,omitempty"`
	TradeOpenedID            string           `json:"tradeOpenedID,omitempty"`
	TradeID                  string           `json:"tradeID,omitempty"`
	TradeReducedID           string           `json:"tradeReducedID,omitempty"`
	TradeClosedIDs           []string         `json:"tradeClosedIDs,omitempty"`
	CancellingTransactionID  string           `json:"cancellingTransactionID,omitempty"`
	CancelledTime            time.Time        `json:"cancelledTime,omitempty"`
	ReplacesOrderID          string           `json:"replacesOrderID,omitempty"`
	ReplacedByOrderID        string           `json:"replacedByOrderID,omitempty"`
	TriggerCondition         string           `json:"triggerCondition,omitempty"`
	GTDTime                  time.Time        `json:"gtdTime,omitempty"`
	Distance                 string           `json:"distance,omitempty"`
}

type OrderPayload struct {
	Order OrderBody `json:"order"`
}

type OrderResponse struct {
	LastTransactionID      string `json:"lastTransactionID"`
	OrderCreateTransaction struct {
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
	} `json:"orderCreateTransaction,omitempty"`
	OrderFillTransaction struct {
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
		Type      string `json:"type"`
		Units     string `json:"units"`
		UserID    int    `json:"userID"`
		FullPrice FullPrice `json:"fullPrice,omitempty"`
		TradesClosed []struct {
			TradeID    string `json:"tradeID"`
			Units      string `json:"units"`
			Financing  string `json:"financing"`
			RealizedPL string `json:"realizedPL"`
			Price      string `json:"price"`
		} `json:"tradesClosed,omitempty"`
		TradeReduced struct {
			TradeID    string `json:"tradeID"`
			Units      string `json:"units"`
			Financing  string `json:"financing"`
			RealizedPL string `json:"realizedPL"`
			Price      string `json:"price"`
		} `json:"tradeReduced,omitempty"`
	} `json:"orderFillTransaction,omitempty"`
	OrderCancelTransaction struct {
		ID     string    `json:"id"`
		Time   time.Time `json:"time"`
		Reason string    `json:"reason"`
	} `json:"orderCancelTransaction,omitempty"`
	RelatedTransactionIDs []string         `json:"relatedTransactionIDs"`
	OrderClientExtensions *OrderExtensions `json:"orderClientExtensions,omitempty"`
	TradeClientExtensions *OrderExtensions `json:"tradeClientExtensions,omitempty"`
}

// GetOrderState returns the state of the order based on the transactions in the response
func (or *OrderResponse) GetOrderState() string {
	if or.OrderCancelTransaction.ID != "" {
		return "CANCELLED"
	}
	if or.OrderFillTransaction.ID != "" {
		return "FILLED"
	}
	if or.OrderCreateTransaction.ID != "" {
		return "PENDING"
	}
	return "UNKNOWN"
}

type OrderInfo struct {
	ID                       string           `json:"id"`
	CreateTime               time.Time        `json:"createTime"`
	State                    string           `json:"state"`
	ClientExtensions         *OrderExtensions `json:"clientExtensions,omitempty"`
	Instrument               string           `json:"instrument"`
	Units                    string           `json:"units"` // Keep as string for API response
	TimeInForce              string           `json:"timeInForce"`
	PriceBound               string           `json:"priceBound,omitempty"`
	Type                     string           `json:"type"`
	PositionFill             string           `json:"positionFill"`
	Price                    string           `json:"price,omitempty"`
	TakeProfitOnFill         *OnFill          `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill           *OnFill          `json:"stopLossOnFill,omitempty"`
	GuaranteedStopLossOnFill *OnFill          `json:"guaranteedStopLossOnFill,omitempty"`
	TrailingStopLossOnFill   *OnFill          `json:"trailingStopLossOnFill,omitempty"`
	TradeClientExtensions    *OrderExtensions `json:"tradeClientExtensions,omitempty"`
	FillingTransactionID     string           `json:"fillingTransactionID,omitempty"`
	FilledTime               time.Time        `json:"filledTime,omitempty"`
	TradeOpenedID            string           `json:"tradeOpenedID,omitempty"`
	TradeID                  string           `json:"tradeID,omitempty"`
	TradeReducedID           string           `json:"tradeReducedID,omitempty"`
	CancellingTransactionID  string           `json:"cancellingTransactionID,omitempty"`
	CancelledTime            time.Time        `json:"cancelledTime,omitempty"`
	ReplacesOrderID          string           `json:"replacesOrderID,omitempty"`
	ReplacedByOrderID        string           `json:"replacedByOrderID,omitempty"`
	TriggerCondition         string           `json:"triggerCondition,omitempty"`
	GTDTime                  time.Time        `json:"gtdTime,omitempty"`
	PartialFill              string           `json:"partialFill,omitempty"`
	Distance                 string           `json:"distance,omitempty"`
}

type RetrievedOrders struct {
	LastTransactionID string      `json:"lastTransactionID"`
	Orders            []OrderInfo `json:"orders,omitempty"`
}

type RetrievedOrder struct {
	Order OrderInfo `json:"order"`
}

type CancelledOrder struct {
	OrderCancelTransaction struct {
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
	} `json:"orderCancelTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
	LastTransactionID     string   `json:"lastTransactionID"`
}

func (c *Connection) CreateOrder(body OrderPayload) (OrderResponse, error) {
	or := OrderResponse{}
	err := c.postAndUnmarshal("/accounts/"+c.accountID+"/orders", body, &or)
	return or, err
}

func (c *Connection) GetOrders(instrument string) (RetrievedOrders, error) {
	endpoint := "/accounts/" + c.accountID + "/orders"
	if instrument != "" {
		endpoint = endpoint + "?instrument=" + instrument
	}

	ro := RetrievedOrders{}
	err := c.getAndUnmarshal(endpoint, &ro)
	return ro, err
}

func (c *Connection) GetPendingOrders() (RetrievedOrders, error) {
	ro := RetrievedOrders{}
	err := c.getAndUnmarshal("/accounts/"+c.accountID+"/pendingOrders", &ro)
	return ro, err
}

func (c *Connection) GetOrder(orderSpecifier string) (RetrievedOrder, error) {
	ro := RetrievedOrder{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/orders/"+
			orderSpecifier,
		&ro,
	)
	return ro, err
}

func (c *Connection) UpdateOrder(orderSpecifier string, body OrderPayload) (RetrievedOrder, error) {
	ro := RetrievedOrder{}
	err := c.putAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/orders/"+
			orderSpecifier,
		body,
		&ro,
	)
	return ro, err
}

func (c *Connection) CancelOrder(orderSpecifier string) (CancelledOrder, error) {
	co := CancelledOrder{}
	err := c.putAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/orders/"+
			orderSpecifier+
			"/cancel",
		nil,
		&co,
	)
	return co, err
}
