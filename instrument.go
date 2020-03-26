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

type OrderBook struct {
	OrderBook   struct {
		Instrument  string    `json:"instrument"`
		Time        time.Time `json:"time"`
		UnixTime    string       `json:"unixTime"`
		Price       string    `json:"price"`
		BucketWidth string    `json:"bucketWidth"`
		Buckets     []Bucket  `json:"buckets"`
	}
}

type PositionBook struct {
	PositionBook   struct {
		Instrument  string    `json:"instrument"`
		Time        time.Time `json:"time"`
		UnixTime    string       `json:"unixTime"`
		Price       string    `json:"price"`
		BucketWidth string    `json:"bucketWidth"`
		Buckets     []Bucket  `json:"buckets"`
	}
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

func (c *OandaConnection) GetCandles(instrument string, count string, granularity string) (InstrumentHistory, error) {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity
	candles, err := c.Request(endpoint)
	if err != nil {
		return InstrumentHistory{}, err
	}
	data := InstrumentHistory{}
	unmarshalJson(candles, &data)

	return data, nil
}

func (c *OandaConnection) GetBidAskCandles(instrument string, count string, granularity string) (BidAskCandles, error) {
	endpoint := "/instruments/" + instrument + "/candles?count=" + count + "&granularity=" + granularity + "&price=BA"
	candles, err := c.Request(endpoint)
	if err != nil {
		return BidAskCandles{}, err
	}
	data := BidAskCandles{}
	unmarshalJson(candles, &data)

	return data, nil
}

func (c *OandaConnection) OrderBook(instrument string) (OrderBook, error) {
	endpoint := "/instruments/" + instrument + "/orderBook"
	orderbook, err := c.Request(endpoint)
	if err != nil {
		return OrderBook{}, err
	}
	data := OrderBook{}
	unmarshalJson(orderbook, &data)

	return data, nil
}

func (c *OandaConnection) PositionBook(instrument string) (PositionBook, error) {
	endpoint := "/instruments/" + instrument + "/positionBook"
	orderbook, err := c.Request(endpoint)
	if err != nil {
		return PositionBook{}, err
	}
	data := PositionBook{}
	unmarshalJson(orderbook, &data)

	return data, nil
}

func (c *OandaConnection) GetInstrumentPrice(instrument string) (InstrumentPricing, error) {
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + instrument
	pricing, err := c.Request(endpoint)
	if err != nil {
		return InstrumentPricing{}, err
	}
	data := InstrumentPricing{}
	unmarshalJson(pricing, &data)

	return data, nil
}
