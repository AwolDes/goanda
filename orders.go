package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/order-ep/

import (
	"encoding/json"
	"time"
)

type OrderExtensions struct {
	Comment string `json:"comment,omitempty"`
	ID      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type OnFill struct {
	TimeInForce string `json:"timeInForce,omitempty"`
	Price       string `json:"price,omitempty"` // must be a string for float precision
}

type OrderBody struct {
	Units            int              `json:"units"`
	Instrument       string           `json:"instrument"`
	TimeInForce      string           `json:"timeInForce"`
	Type             string           `json:"type"`
	PositionFill     string           `json:"positionFill,omitempty"`
	Price            string           `json:"price,omitempty"`
	TakeProfitOnFill *OnFill          `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill   *OnFill          `json:"stopLossOnFill,omitempty"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
	TradeID          string           `json:"tradeId,omitempty"`
}

type OrderPayload struct {
	Order OrderBody `json:"order"`
}
type OrderResponse struct {
	LastTransactionID      string `json:"lastTransactionID"`
	OrderCreateTransaction struct {
		AccountID    string    `json:"accountID"`
		BatchID      string    `json:"batchID"`
		ID           string    `json:"id"`
		Instrument   string    `json:"instrument"`
		PositionFill string    `json:"positionFill"`
		Reason       string    `json:"reason"`
		Time         time.Time `json:"time"`
		TimeInForce  string    `json:"timeInForce"`
		Type         string    `json:"type"`
		Units        string    `json:"units"`
		UserID       int       `json:"userID"`
	} `json:"orderCreateTransaction"`
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
		Type   string `json:"type"`
		Units  string `json:"units"`
		UserID int    `json:"userID"`
	} `json:"orderFillTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
}

type OrderInfo struct {
	ClientExtensions struct {
		Comment string `json:"comment,omitempty"`
		ID      string `json:"id,omitempty"`
		Tag     string `json:"tag,omitempty"`
	} `json:"clientExtensions,omitempty"`
	CreateTime       time.Time `json:"createTime"`
	ID               string    `json:"id"`
	Instrument       string    `json:"instrument,omitempty"`
	PartialFill      string    `json:"partialFill"`
	PositionFill     string    `json:"positionFill"`
	Price            string    `json:"price"`
	ReplacesOrderID  string    `json:"replacesOrderID,omitempty"`
	State            string    `json:"state"`
	TimeInForce      string    `json:"timeInForce"`
	TriggerCondition string    `json:"triggerCondition"`
	Type             string    `json:"type"`
	Units            string    `json:"units,omitempty"`
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

func (c *OandaConnection) CreateOrder(body OrderPayload) OrderResponse {
	endpoint := "/accounts/" + c.accountID + "/orders"
	jsonBody, err := json.Marshal(body)
	checkErr(err)
	response := c.Send(endpoint, jsonBody)
	data := OrderResponse{}
	unmarshalJson(response, &data)
	return data
}

func (c *OandaConnection) GetOrders(instrument string) RetrievedOrders {
	endpoint := "/accounts/" + c.accountID + "/orders"

	if instrument != "" {
		endpoint = endpoint + "?instrument=" + instrument
	}

	response := c.Request(endpoint)
	data := RetrievedOrders{}
	unmarshalJson(response, &data)
	return data
}

func (c *OandaConnection) GetPendingOrders() RetrievedOrders {
	endpoint := "/accounts/" + c.accountID + "/pendingOrders"

	response := c.Request(endpoint)
	data := RetrievedOrders{}
	unmarshalJson(response, &data)
	return data

}

func (c *OandaConnection) GetOrder(orderSpecifier string) RetrievedOrder {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier

	response := c.Request(endpoint)
	data := RetrievedOrder{}
	unmarshalJson(response, &data)
	return data

}

func (c *OandaConnection) UpdateOrder(orderSpecifier string, body OrderPayload) RetrievedOrder {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier
	jsonBody, err := json.Marshal(body)
	checkErr(err)
	response := c.Update(endpoint, jsonBody)
	data := RetrievedOrder{}
	unmarshalJson(response, &data)
	return data
}

func (c *OandaConnection) CancelOrder(orderSpecifier string) CancelledOrder {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier + "/cancel"
	response := c.Update(endpoint, nil)
	data := CancelledOrder{}
	unmarshalJson(response, &data)
	return data

}
