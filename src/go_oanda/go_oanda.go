package go_oanda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Headers struct {
	contentType     string
	agent           string
	datetime_format string
	auth            string
}

type OandaConnection struct {
	hostname        string
	port            int
	ssl             bool
	token           string
	datetime_format string
	headers         *Headers
}

const OANDA_AGENT string = "v20-golang/0.0.1"

func NewConnection(token string, live bool) *OandaConnection {
	hostname := ""
	// should we use the live API?
	if live {
		hostname = "https://api-fxtrade.oanda.com/"
	} else {
		hostname = "https://api-fxpractice.oanda.com/"
	}

	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString("Bearer ")
	buffer.WriteString(token)

	authHeader := buffer.String()
	// Create headers for oanda to be used in requests
	headers := &Headers{
		contentType:     "application/json",
		agent:           OANDA_AGENT,
		datetime_format: "RFC3339",
		auth:            authHeader,
	}
	// Create the connection object
	connection := &OandaConnection{
		hostname: hostname,
		port:     443,
		ssl:      true,
		token:    token,
		headers:  headers,
	}

	return connection
}

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

func (c *OandaConnection) Request(endpoint string) InstrumentHistory {
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}

	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(c.hostname)
	buffer.WriteString(endpoint)

	url := buffer.String()

	// New request object
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)

	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	data := InstrumentHistory{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println(data)
	return data
}
