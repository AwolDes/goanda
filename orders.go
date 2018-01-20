package goanda

type Order struct {
	Units string `json:"units"`
	Instrument string `json:"instrument"`
	TimeInForce string `json:timeInForce`
	Type string `json:"type"`
	PositionFill string `json:"positionFill"`
}

// func (c *OandaConnection) CreateOrder(body Order) {
// 	endpoint := "/accounts/" + c.accountId + "/orders"

// 	response := c.Send(endpoint, body)
// 	data := InstrumentHistory{}
// 	unmarshalJson(response, &data)

// 	return data
// }