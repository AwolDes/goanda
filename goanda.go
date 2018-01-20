package goanda

import (
	"bytes"
	"io/ioutil"
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
		hostname = "https://api-fxtrade.oanda.com/v3"
	} else {
		hostname = "https://api-fxpractice.oanda.com/v3"
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

func (c *OandaConnection) Request(endpoint string) []byte {
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
	checkErr(err)

	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)

	res, getErr := client.Do(req)
	checkErr(getErr)

	body, readErr := ioutil.ReadAll(res.Body)
	checkErr(readErr)

	return body
}
