package goanda

import (
	"net/url"
	"time"
)

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/transaction-ep/

type TransactionPages struct {
	Count             int       `json:"count"`
	From              time.Time `json:"from"`
	LastTransactionID string    `json:"lastTransactionID"`
	PageSize          int       `json:"pageSize"`
	Pages             []string  `json:"pages"`
	To                time.Time `json:"to"`
}

type Transaction struct {
	LastTransactionID string `json:"lastTransactionID"`
	Transaction       struct {
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
	} `json:"transaction"`
}

type Transactions struct {
	LastTransactionID string `json:"lastTransactionID"`
	Transactions      []struct {
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
	} `json:"transactions"`
}

// https://golang.org/pkg/time/#Time.AddDate
// https://play.golang.org/p/Dw7D4JJ7EC
func (c *Connection) GetTransactions(from time.Time, to time.Time) (TransactionPages, error) {
	tp := TransactionPages{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/transactions?to="+
			url.QueryEscape(to.Format(time.RFC3339))+
			"&from="+
			url.QueryEscape(from.Format(time.RFC3339)),
		&tp,
	)
	return tp, err
}

func (c *Connection) GetTransaction(ticket string) (Transaction, error) {
	tr := Transaction{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/transactions/"+
			ticket,
		&tr,
	)
	return tr, err
}

func (c *Connection) GetTransactionsSinceId(id string) (Transactions, error) {
	tr := Transactions{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/transactions/sinceid?id="+
			id,
		&tr,
	)
	return tr, err
}
