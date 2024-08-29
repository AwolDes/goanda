package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/trade-ep/

import (
	"time"
)

type ReceivedTrades struct {
	LastTransactionID string  `json:"lastTransactionID"`
	Trades            []Trade `json:"trades"`
}

type ReceivedTrade struct {
	LastTransactionID string `json:"lastTransactionID"`
	Trade             Trade  `json:"trade"`
}

type Trade struct {
	ID                    string                 `json:"id"`
	Instrument            string                 `json:"instrument"`
	Price                 string                 `json:"price"`
	OpenTime              time.Time              `json:"openTime"`
	State                 string                 `json:"state"`
	InitialUnits          string                 `json:"initialUnits"`
	InitialMarginRequired string                 `json:"initialMarginRequired"`
	CurrentUnits          string                 `json:"currentUnits"`
	RealizedPL            string                 `json:"realizedPL"`
	UnrealizedPL          string                 `json:"unrealizedPL"`
	MarginUsed            string                 `json:"marginUsed"`
	AverageClosePrice     string                 `json:"averageClosePrice,omitempty"`
	ClosingTransactionIDs []string               `json:"closingTransactionIDs,omitempty"`
	Financing             string                 `json:"financing"`
	CloseTime             time.Time              `json:"closeTime,omitempty"`
	ClientExtensions      *OrderExtensions       `json:"clientExtensions,omitempty"`
	TakeProfitOrder       *TakeProfitOrder       `json:"takeProfitOrder,omitempty"`
	StopLossOrder         *StopLossOrder         `json:"stopLossOrder,omitempty"`
	TrailingStopLossOrder *TrailingStopLossOrder `json:"trailingStopLossOrder,omitempty"`
}

type TakeProfitOrder struct {
	ID               string           `json:"id"`
	CreateTime       time.Time        `json:"createTime"`
	Type             string           `json:"type"`
	TradeID          string           `json:"tradeID"`
	ClientTradeID    string           `json:"clientTradeID,omitempty"`
	Price            string           `json:"price"`
	TimeInForce      string           `json:"timeInForce"`
	TriggerCondition string           `json:"triggerCondition"`
	State            string           `json:"state"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
}

type StopLossOrder struct {
	ID               string           `json:"id"`
	CreateTime       time.Time        `json:"createTime"`
	Type             string           `json:"type"`
	TradeID          string           `json:"tradeID"`
	ClientTradeID    string           `json:"clientTradeID,omitempty"`
	Price            string           `json:"price"`
	TimeInForce      string           `json:"timeInForce"`
	TriggerCondition string           `json:"triggerCondition"`
	State            string           `json:"state"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
	Guaranteed       bool             `json:"guaranteed"`
}

type TrailingStopLossOrder struct {
	ID               string           `json:"id"`
	CreateTime       time.Time        `json:"createTime"`
	Type             string           `json:"type"`
	TradeID          string           `json:"tradeID"`
	ClientTradeID    string           `json:"clientTradeID,omitempty"`
	Distance         string           `json:"distance"`
	TimeInForce      string           `json:"timeInForce"`
	TriggerCondition string           `json:"triggerCondition"`
	State            string           `json:"state"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
}

type CloseTradePayload struct {
	Units string `json:"units"`
}

type ModifiedTrade struct {
	OrderCreateTransaction struct {
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
	} `json:"orderCreateTransaction"`
	OrderFillTransaction struct {
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
	} `json:"orderFillTransaction"`
	OrderCancelTransaction struct {
		Type      string    `json:"type"`
		OrderID   string    `json:"orderID"`
		Reason    string    `json:"reason"`
		ID        string    `json:"id"`
		UserID    int       `json:"userID"`
		AccountID string    `json:"accountID"`
		BatchID   string    `json:"batchID"`
		RequestID string    `json:"requestID"`
		Time      time.Time `json:"time"`
	} `json:"orderCancelTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
	LastTransactionID     string   `json:"lastTransactionID"`
}

type FullPrice struct {
	CloseoutBid string       `json:"closeoutBid"`
	CloseoutAsk string       `json:"closeoutAsk"`
	Timestamp   string       `json:"timestamp"`
	Bids        []PriceLevel `json:"bids"`
	Asks        []PriceLevel `json:"asks"`
}

type PriceLevel struct {
	Price     string `json:"price"`
	Liquidity string `json:"liquidity"`
}

func (c *Connection) GetTradesForInstrument(instrument string) (ReceivedTrades, error) {
	rt := ReceivedTrades{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/trades"+
			"?instrument="+
			instrument,
		&rt,
	)
	return rt, err
}

func (c *Connection) GetOpenTrades() (ReceivedTrades, error) {
	rt := ReceivedTrades{}
	err := c.getAndUnmarshal("/accounts/"+c.accountID+"/openTrades", &rt)
	return rt, err
}

func (c *Connection) GetTrade(ticket string) (ReceivedTrade, error) {
	rt := ReceivedTrade{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/trades/"+
			ticket,
		&rt,
	)
	return rt, err
}

// Default is close the whole position using the string "ALL" in body.units
func (c *Connection) ReduceTradeSize(ticket string, body CloseTradePayload) (ModifiedTrade, error) {
	mt := ModifiedTrade{}
	err := c.putAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/trades/"+
			ticket+
			"/close",
		body,
		&mt,
	)
	return mt, err
}
