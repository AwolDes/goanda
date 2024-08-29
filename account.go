package goanda

import (
	"time"
)

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/account-ep/

type AccountProperties struct {
	ID           string `json:"id"`
	Mt4AccountID int    `json:"mt4AccountID"`
	Tags         []string
}

type AccountInfo struct {
	Account struct {
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
	} `json:"account"`
	LastTransactionID string `json:"lastTransactionID"`
}

type AccountSummary struct {
	Account struct {
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
	} `json:"account"`
	LastTransactionID string `json:"lastTransactionID"`
}

type Instruments []struct {
	DisplayName                 string `json:"displayName"`
	DisplayPrecision            int    `json:"displayPrecision"`
	MarginRate                  string `json:"marginRate"`
	MaximumOrderUnits           string `json:"maximumOrderUnits"`
	MaximumPositionSize         string `json:"maximumPositionSize"`
	MaximumTrailingStopDistance string `json:"maximumTrailingStopDistance"`
	MinimumTradeSize            string `json:"minimumTradeSize"`
	MinimumTrailingStopDistance string `json:"minimumTrailingStopDistance"`
	Name                        string `json:"name"`
	PipLocation                 int    `json:"pipLocation"`
	TradeUnitsPrecision         int    `json:"tradeUnitsPrecision"`
	Type                        string `json:"type"`
}

type AccountChanges struct {
	Changes struct {
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
	} `json:"changes"`
	LastTransactionID string `json:"lastTransactionID"`
	State             struct {
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
	} `json:"state"`
}

type OrderDetails struct {
	GainPerPipPerMillionUnits float64 `json:"gainPerPipPerMillionUnits,string"`
	LossPerPipPerMillionUnits float64 `json:"lossPerPipPerMillionUnits,string"`
	UnitsAvailable            struct {
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
	} `json:"unitsAvailable"`
	UnitValues struct {
		Isolation struct {
			Units               float64 `json:"units,string"`
			Commission          float64 `json:"commission,string"`
			PositionValueChange float64 `json:"positionValueChange,string"`
			PositionValue       float64 `json:"positionValue,string"`
			MarginRequired      float64 `json:"marginRequired,string"`
			MarginUsed          float64 `json:"marginUsed,string"`
		} `json:"isolation"`
		PositionDefault struct {
			Units               float64 `json:"units,string"`
			Commission          float64 `json:"commission,string"`
			PositionValueChange float64 `json:"positionValueChange,string"`
			PositionValue       float64 `json:"positionValue,string"`
			MarginRequired      float64 `json:"marginRequired,string"`
			MarginUsed          float64 `json:"marginUsed,string"`
		} `json:"positionDefault"`
		PositionOpenOnly struct {
			Units               float64 `json:"units,string"`
			Commission          float64 `json:"commission,string"`
			PositionValueChange float64 `json:"positionValueChange,string"`
			PositionValue       float64 `json:"positionValue,string"`
			MarginRequired      float64 `json:"marginRequired,string"`
			MarginUsed          float64 `json:"marginUsed,string"`
		} `json:"positionOpenOnly"`
		PositionReduceFirst struct {
			Units               float64 `json:"units,string"`
			Commission          float64 `json:"commission,string"`
			PositionValueChange float64 `json:"positionValueChange,string"`
			PositionValue       float64 `json:"positionValue,string"`
			MarginRequired      float64 `json:"marginRequired,string"`
			MarginUsed          float64 `json:"marginUsed,string"`
		} `json:"positionReduceFirst"`
		PositionReduceOnly struct {
			Units               float64 `json:"units,string"`
			Commission          float64 `json:"commission,string"`
			PositionValueChange float64 `json:"positionValueChange,string"`
			PositionValue       float64 `json:"positionValue,string"`
			MarginRequired      float64 `json:"marginRequired,string"`
			MarginUsed          float64 `json:"marginUsed,string"`
		} `json:"positionReduceOnly"`
	} `json:"unitValues"`
	ValueTables struct {
		CommissionTable []struct {
			Units string `json:"units"`
			Value string `json:"value"`
		} `json:"commissionTable"`
	} `json:"valueTables"`
	LastTransactionID string `json:"lastTransactionID"`
}

// Accounts returns a slice of information on accounts authorized for the token.
func (c *Connection) Accounts() ([]AccountProperties, error) {
	var ap []AccountProperties
	err := c.getAndUnmarshal("/accounts", &ap)
	return ap, err
}

// GetAccount returns information on the account.
func (c *Connection) GetAccount(id string) (AccountInfo, error) {
	ai := AccountInfo{}
	err := c.getAndUnmarshal("/accounts/"+id, &ai)
	return ai, err
}

func (c *Connection) GetOrderDetails(instrument string, units string) (OrderDetails, error) {
	od := OrderDetails{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/orderEntryData?disableFiltering=true&instrument="+
			instrument+
			"&orderPositionFill=DEFAULT&units="+
			units,
		&od,
	)
	return od, err
}

func (c *Connection) GetAccountSummary() (AccountSummary, error) {
	as := AccountSummary{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/summary",
		&as,
	)
	return as, err
}

func (c *Connection) GetAccountInstruments(id string) (Instruments, error) {
	var response struct {
		Instruments Instruments `json:"instruments"`
	}

	err := c.getAndUnmarshal(
		"/accounts/"+
			id+
			"/instruments",
		&response,
	)

	return response.Instruments, err
}

func (c *Connection) GetAccountChanges(id string, transactionId string) (AccountChanges, error) {
	ac := AccountChanges{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			id+
			"/changes?sinceTransactionID="+
			transactionId,
		&ac,
	)
	return ac, err
}
