package goanda

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
	Price       string `json:"price,omitempty"`
}

type OrderBody struct {
	Units            string          `json:"units"`
	Instrument       string          `json:"instrument"`
	TimeInForce      string          `json:"timeInForce"`
	Type             string          `json:"type"`
	PositionFill     string          `json:"positionFill"`
	Price            string          `json:"price,omitempty"`
	StopLossOnFill   *OnFill         `json:"stopLossOnFill,omitempty"`
	ClientExtensions OrderExtensions `json:"clientExtensions,omitempty"`
	TradeID          string          `json:"tradeId,omitempty"`
}

type Order struct {
	Order OrderBody `json:"order"`
}

type TransactionInfo struct {
	AccountID    string    `json:"accountID"`
	BatchID      string    `json:"batchID"`
	ID           string    `json:"id"`
	Instrument   string    `json:"instrument"`
	PositionFill string    `json:"positionFill"`
	Reason       string    `json:"reason"`
	Time         time.Time `json:"time"`
	TimeInForce  string    `json:"timeInForce"`
	OrderType    string    `json:"OrderType"`
	Units        string    `json:"units"`
	UserID       int       `json:"userID"`
}
type TradeOpen struct {
	TradeID string `json:"tradeId"`
	Units   string `json:"units"`
}

type TransactionAccountInfo struct {
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
	TradeOpened    TradeOpen `json:""`
	OrderType      string    `json:"OrderType"`
	Units          string    `json:"units"`
	UserID         int       `json:"userID"`
}

type OrderResponse struct {
	LastTransactionID      string                 `json:"lastTransactionID"`
	OrderCreateTransaction TransactionInfo        `json:"orderCreateTransaction"`
	OrderFillTransaction   TransactionAccountInfo `json:"orderFillTransaction"`
	RelatedTransactionIDs  []string               `json:"relatedTransactionIDs"`
}

type RetrievedOrder struct {
	ClientExtensions struct {
		Comment string `json:"comment"`
		ID      string `json:"id"`
		Tag     string `json:"tag"`
	} `json:"clientExtensions"`
	CreateTime       time.Time `json:"createTime"`
	ID               string    `json:"id"`
	Instrument       string    `json:"instrument"`
	PartialFill      string    `json:"partialFill"`
	PositionFill     string    `json:"positionFill"`
	Price            string    `json:"price"`
	ReplacesOrderID  string    `json:"replacesOrderID"`
	State            string    `json:"state"`
	TimeInForce      string    `json:"timeInForce"`
	TriggerCondition string    `json:"triggerCondition"`
	Type             string    `json:"type"`
	Units            string    `json:"units"`
}

type RetreievedOrders struct {
	LastTransactionID string           `json:"lastTransactionID"`
	Orders            []RetrievedOrder `json:"orders"`
}

func (c *OandaConnection) CreateOrder(body *Order) OrderResponse {
	endpoint := "/accounts/" + c.accountID + "/orders"
	jsonBody, err := json.Marshal(body)
	checkErr(err)
	response := c.Send(endpoint, jsonBody)
	data := OrderResponse{}
	unmarshalJson(response, &data)
	return data
}

func (c *OandaConnection) GetOrders(instrument string) RetreievedOrders {
	endpoint := "/accounts/" + c.accountID + "/orders"

	if instrument != "" {
		endpoint = endpoint + "?instrument=" + instrument
	}

	response := c.Request(endpoint)
	data := RetreievedOrders{}
	unmarshalJson(response, &data)
	return data

}
