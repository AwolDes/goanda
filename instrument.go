package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/instrument-ep/

import (
	"errors"
	"strconv"
	"time"
)

// GranularityFromDuration tries to find a granularity for the given duration
func GranularityFromDuration(d time.Duration) (Granularity, error) {
	if _, ok := candlestickGranularity[Granularity(d)]; ok {
		return Granularity(d), nil
	}
	return 0, errors.New("No such granularity")
}

// Granularity defines a candle's time period
type Granularity time.Duration

// Duration returns the granularity as a time.Duration
func (g Granularity) Duration() time.Duration {
	return time.Duration(g)
}

// String returns the granularity as a string, formatted to the oanda standard
func (g Granularity) String() string {
	return candlestickGranularity[g]
}

// Granularities available to the API
const (
	GranularityFiveSeconds    = Granularity(time.Second * 5)
	GranularityTenSeconds     = Granularity(time.Second * 10)
	GranularityFifteenSeconds = Granularity(time.Second * 15)
	GranularityThirtySeconds  = Granularity(time.Second * 30)
	GranularityMinute         = Granularity(time.Minute)
	GranularityTwoMinutes     = Granularity(time.Minute * 2)
	GranularityFourMinutes    = Granularity(time.Minute * 4)
	GranularityFiveMinutes    = Granularity(time.Minute * 5)
	GranularityTenMinutes     = Granularity(time.Minute * 10)
	GranularityFifteenMinutes = Granularity(time.Minute * 15)
	GranularityThirtyMinutes  = Granularity(time.Minute * 30)
	GranularityHour           = Granularity(time.Hour)
	GranularityTwoHours       = Granularity(time.Hour * 2)
	GranularityThreeHours     = Granularity(time.Hour * 3)
	GranularityFourHours      = Granularity(time.Hour * 4)
	GranularitySixHours       = Granularity(time.Hour * 6)
	GranularityEightHours     = Granularity(time.Hour * 8)
	GranularityTwelveHours    = Granularity(time.Hour * 12)
	GranularityDay            = Granularity(time.Hour * 24)
	GranularityWeek           = Granularity(time.Hour * 24 * 7)
	GranularityMonth          = Granularity(time.Hour * 24 * 7 * 30)
)

var candlestickGranularity = map[Granularity]string{
	GranularityFiveSeconds:    "S5",
	GranularityTenSeconds:     "S10",
	GranularityFifteenSeconds: "S15",
	GranularityThirtySeconds:  "S30",
	GranularityMinute:         "M1",
	GranularityTwoMinutes:     "M2",
	GranularityFourMinutes:    "M4",
	GranularityFiveMinutes:    "M5",
	GranularityTenMinutes:     "M10",
	GranularityFifteenMinutes: "M15",
	GranularityThirtyMinutes:  "M30",
	GranularityHour:           "H1",
	GranularityTwoHours:       "H2",
	GranularityThreeHours:     "H3",
	GranularityFourHours:      "H4",
	GranularitySixHours:       "H6",
	GranularityEightHours:     "H8",
	GranularityTwelveHours:    "H12",
	GranularityDay:            "D",
	GranularityWeek:           "W",
	GranularityMonth:          "M",
}

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

func (c *Connection) GetCandles(instrument string, count int, g Granularity) (InstrumentHistory, error) {
	ca := InstrumentHistory{}
	err := c.getAndUnmarshal(
		"/instruments/"+
			instrument+
			"/candles?count="+
			strconv.Itoa(count)+
			"&granularity="+
			g.String(),
		&ca,
	)
	return ca, err
}

func (c *Connection) GetTimeToCandles(instrument string, count int, g Granularity, to time.Time) (InstrumentHistory, error) {
	ih := InstrumentHistory{}
	err := c.requestAndUnmarshal(
		"/instruments/"+
			instrument+
			"/candles?count="+
			strconv.Itoa(count)+
			"&to="+
			strconv.Itoa(int(to.Unix()))+
			"&granularity="+
			g.String(),
		&ih,
	)
	return ih, err
}
func (c *Connection) GetTimeFromCandles(instrument string, count int, g Granularity, from time.Time) (InstrumentHistory, error) {
	ih := InstrumentHistory{}
	err := c.requestAndUnmarshal(
		"/instruments/"+
			instrument+
			"/candles?count="+
			strconv.Itoa(count)+
			"&from="+
			strconv.Itoa(int(from.Unix()))+
			"&granularity="+
			g.String(),
		&ih,
	)
	return ih, err
}

func (c *Connection) GetBidAskCandles(instrument string, count string, g Granularity) (BidAskCandles, error) {
	ca := BidAskCandles{}
	err := c.getAndUnmarshal(
		"/instruments/"+
			instrument+
			"/candles?count="+
			count+
			"&granularity="+
			g.String()+
			"&price=BA",
		&ca,
	)
	return ca, err
}

func (c *Connection) OrderBook(instrument string) (BrokerBook, error) {
	bb := BrokerBook{}
	err := c.getAndUnmarshal(
		"/instruments/"+
			instrument+
			"/orderBook",
		&bb,
	)
	return bb, err
}

func (c *Connection) PositionBook(instrument string) (BrokerBook, error) {
	bb := BrokerBook{}
	err := c.getAndUnmarshal(
		"/instruments/"+
			instrument+
			"/positionBook",
		&bb,
	)
	return bb, err
}

func (c *Connection) GetInstrumentPrice(instrument string) (InstrumentPricing, error) {
	ip := InstrumentPricing{}
	err := c.getAndUnmarshal(
		"/accounts/"+
			c.accountID+
			"/pricing?instruments="+
			instrument,
		&ip,
	)
	return ip, err
}
