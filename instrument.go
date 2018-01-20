package goanda

import (
	"encoding/json"
)

type Candle struct {
	Open  string `json:"o"`
	Close string `json:"c"`
	Low   string `json:"l"`
	High  string `json:"h"`
}

type Candles struct {
	Complete bool    `json:"complete"`
	Volume   float64 `json:"volume"`
	Time     string  `json:"time"`
	Mid      Candle  `json:"mid"`
}

type InstrumentHistory struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}

func (c *OandaConnection) GetCandles(instrument string) InstrumentHistory {
	endpoint := "/instruments/" + instrument + "/candles"
	candles := c.Request(endpoint)
	data := InstrumentHistory{}
	jsonErr := json.Unmarshal(candles, &data)
	checkErr(jsonErr)

	return data
}