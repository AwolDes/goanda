package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/position-ep/

type OpenPositions struct {
	LastTransactionID string `json:"lastTransactionID"`
	Positions         []struct {
		Instrument string `json:"instrument"`
		Long       struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"long"`
		Pl           string `json:"pl"`
		ResettablePL string `json:"resettablePL"`
		Short        struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"short"`
		UnrealizedPL string `json:"unrealizedPL"`
	} `json:"positions"`
}

type ClosePositionPayload struct {
	LongUnits  string `json:"longUnits"`
	ShortUnits string `json:"shortUnits"`
}

func (c *Connection) GetOpenPositions() (OpenPositions, error) {
	op := OpenPositions{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/openPositions",
		&op,
	)
	return op, err
}

func (c *Connection) ClosePosition(instrument string, body ClosePositionPayload) (ModifiedTrade, error) {
	mt := ModifiedTrade{}
	err := c.putAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/positions/"+
			instrument+
			"/close",
		body,
		&mt,
	)
	return mt, err
}
