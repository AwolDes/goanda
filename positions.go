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
			Pl           string `json:"pl"`
			ResettablePL string `json:"resettablePL"`
			Units        string `json:"units"`
			UnrealizedPL string `json:"unrealizedPL"`
		} `json:"short"`
		UnrealizedPL string `json:"unrealizedPL"`
	} `json:"positions"`
}

type ClosePositionPayload struct {
	LongUnits  string `json:"longUnits"`
	ShortUnits string `json:"shortUnits"`
}

func (c *OandaConnection) GetOpenPositions() (OpenPositions, error) {
	endpoint := "/accounts/" + c.accountID + "/openPositions"

	response, err := c.Request(endpoint)
	if err != nil {
		return OpenPositions{}, err
	}
	data := OpenPositions{}
	unmarshalJson(response, &data)
	return data, nil
}

func (c *OandaConnection) ClosePosition(instrument string, body ClosePositionPayload) (ModifiedTrade, error) {
	endpoint := "/accounts/" + c.accountID + "/positions/" + instrument + "/close"
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return ModifiedTrade{}, err
	}
	response, err := c.Update(endpoint, jsonBody)
	if err != nil {
		return ModifiedTrade{}, err
	}
	data := ModifiedTrade{}
	unmarshalJson(response, &data)
	return data, nil
}
