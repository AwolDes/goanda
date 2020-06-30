package goanda

import "encoding/json"

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
	ShortUnits string `json: "shortUnits"`
}

func (c *OandaConnection) GetOpenPositions() OpenPositions {
	endpoint := "/accounts/" + c.accountID + "/openPositions"

	response := c.Request(endpoint)
	data := OpenPositions{}
	unmarshalJson(response, &data)
	return data
}

func (c *OandaConnection) ClosePosition(instrument string, body ClosePositionPayload) ModifiedTrade {
	endpoint := "/accounts/" + c.accountID + "/positions/" + instrument + "/close"
	jsonBody, err := json.Marshal(body)
	checkErr(err)
	response := c.Update(endpoint, jsonBody)
	data := ModifiedTrade{}
	unmarshalJson(response, &data)
	return data
}
