package goanda

import (
	"encoding/json"
)

type Order struct {
	Units int `json:"units"`
	Instrument string `json:"instrument"`
	TimeInForce string `json:timeInForce`
	Type string `json:"type"`
	PositionFill string `json:"positionFill"`
}

type TransactionInfo struct {
	AccountID string `json:"accountID"`
	BatchID string `json:"batchID"` 
	Id string `json:"id"`
	Instrument string `json:"instrument"` 
	PositionFill string `json:"positionFill"`
	Reason string `json:"reason"`
	Time string `json:"time"`
	TimeInForce string `json:"timeInForce"`
	OrderType string `json:"OrderType"`
	Units string `json:"units"`
	UserID string `json:"userID"`
}
type TradeOpen struct {
	TradeID string `json:"tradeId"`
	Units string `json:"units"`
}

type TransactionAccountInfo struct {
	AccountBalance string `json:"accountBalance"`
	AccountID string `json:"accountID"`
	BatchID string `json:"batchID"`
	Financing string `json:"financing"`
	Id string `json:"id"`
	Instrument string `json:"instrument"`
	OrderID string `json:"orderID"`
	Pl string `json:"pl"`
	Price string `json:"price"`
	Reason string `json:"reason"`
	Time string `json:"time"`
	TradeOpened TradeOpen `json:""`
	OrderType string `json:"OrderType"`
	Units string `json:"units"`
	UserID string `json:"userID"`
}

type OrderResponse struct {
	LastTransactionID string `json:"lastTransactionID"` 
	OrderCreateTransaction TransactionInfo `json:"orderCreateTransaction"` 
	OrderFillTransaction TransactionAccountInfo `json:"orderFillTransaction"` 
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
}

func (c *OandaConnection) CreateOrder(body *Order) OrderResponse {
	endpoint := "/accounts/" + c.accountId + "/orders"
	jsonBody, err := json.Marshal(body)
	checkErr(err)
	response := c.Send(endpoint, jsonBody)
	data := OrderResponse{}
	unmarshalJson(response, &data)

	return data
}