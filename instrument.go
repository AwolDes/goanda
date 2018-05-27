package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/instrument-ep/

import (
	"time"
)

type Candle struct {
	Open  float64 `json:"o,string"`
	Close float64 `json:"c,string"`
	Low   float64 `json:"l,string"`
	High  float64 `json:"h,string"`
}

type Candles struct {
	Complete bool      `json:"complete"`
	Volume   int       `json:"volume"`
	Time     time.Time `json:"time"`
	Mid      Candle    `json:"mid"`
}

type BidAskCandles struct {
	Candles []struct {
		Ask struct {
			C float64 `json:"c,string"`
			H float64 `json:"h,string"`
			L float64 `json:"l,string"`
			O float64 `json:"o,string"`
		} `json:"ask"`
		Bid struct {
			C float64 `json:"c,string"`
			H float64 `json:"h,string"`
			L float64 `json:"l,string"`
			O float64 `json:"o,string"`
		} `json:"bid"`
		Complete bool      `json:"complete"`
		Time     time.Time `json:"time"`
		Volume   int       `json:"volume"`
	} `json:"candles"`
}

type InstrumentHistory struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}

type Bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}

type BrokerBook struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []Bucket  `json:"buckets"`
}

type InstrumentPricing struct {
	Time   time.Time `json:"time"`
	Prices []struct {
		Type string    `json:"type"`
		Time time.Time `json:"time"`
		Bids []struct {
			Price     float64 `json:"price,string"`
			Liquidity int     `json:"liquidity"`
		} `json:"bids"`
		Asks []struct {
			Price     float64 `json:"price,string"`
			Liquidity int     `json:"liquidity"`
		} `json:"asks"`
		CloseoutBid    float64 `json:"closeoutBid,string"`
		CloseoutAsk    float64 `json:"closeoutAsk,string"`
		Status         string  `json:"status"`
		Tradeable      bool    `json:"tradeable"`
		UnitsAvailable struct {
			Default struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"default"`
			OpenOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"openOnly"`
			ReduceFirst struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceFirst"`
			ReduceOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceOnly"`
		} `json:"unitsAvailable"`
		QuoteHomeConversionFactors struct {
			PositiveUnits string `json:"positiveUnits"`
			NegativeUnits string `json:"negativeUnits"`
		} `json:"quoteHomeConversionFactors"`
		Instrument string `json:"instrument"`
	} `json:"prices"`
}

func (c *OandaConnection) GetCandles(instrument string, count string, granularity string) InstrumentHistory {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity
	candles := c.Request(endpoint)
	data := InstrumentHistory{}
	unmarshalJson(candles, &data)

	return data
}

func (c *OandaConnection) GetBidAskCandles(instrument string, count string, granularity string) BidAskCandles {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity + "&price=BA"
	candles := c.Request(endpoint)
	data := BidAskCandles{}
	unmarshalJson(candles, &data)

	return data
}

func (c *OandaConnection) OrderBook(instrument string) BrokerBook {
	endpoint := "/instruments/" + instrument + "/orderBook"
	orderbook := c.Request(endpoint)
	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data
}

func (c *OandaConnection) PositionBook(instrument string) BrokerBook {
	endpoint := "/instruments/" + instrument + "/positionBook"
	orderbook := c.Request(endpoint)
	data := BrokerBook{}
	unmarshalJson(orderbook, &data)

	return data
}

func (c *OandaConnection) GetInstrumentPrice(instrument string) InstrumentPricing {
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + instrument
	pricing := c.Request(endpoint)
	data := InstrumentPricing{}
	unmarshalJson(pricing, &data)

	return data
}
