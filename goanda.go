package goanda

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiUserAgent = "v20-golang/0.0.1"
	httpTimeout  = time.Second * 5
)

// ConnectionConfig is used to configure new connections
// Defaults;
//	UserAgent	= v20-golang/0.0.1
//	Timeout		= 5 seconds
//	Live		= False
type ConnectionConfig struct {
	UserAgent string
	Timeout   time.Duration
	Live      bool
}

// Connection describes a connection to the Oanda v20 API
// It is thread safe
type Connection struct {
	hostname   string
	accountID  string
	authHeader string
	userAgent  string
	client     http.Client
}

// NewConnection creates a new connection
// This function calls Connection.CheckConnection(), returning any errors
// Supplying a config is optional, with sane defaults (paper trading) being used otherwise.
func NewConnection(accountID string, token string, config *ConnectionConfig) (*Connection, error) {
	// Make new connection with defaults
	nc := &Connection{
		hostname:   "https://api-fxpractice.oanda.com/v3",
		accountID:  accountID,
		authHeader: "Bearer " + token,
		userAgent:  apiUserAgent,
		client: http.Client{
			Timeout: httpTimeout,
		},
	}

	// Overwrite things if we've been given configuration for them
	if config != nil {
		if config.Live {
			nc.hostname = "https://api-fxtrade.oanda.com/v3"
		}

		if config.Timeout != 0 {
			nc.client = http.Client{
				Timeout: config.Timeout,
			}
		}

		if config.UserAgent != "" {
			nc.userAgent = config.UserAgent
		}
	}

	return nc, nc.CheckConnection()
}

// CheckConnection performs a request, returning any errors encountered
func (c *Connection) CheckConnection() error {
	_, err := c.Get("/accounts/" + c.accountID)
	return err
}

// Get performs a genereic http get on the api
func (c *Connection) Get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.hostname+endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.makeRequest(endpoint, c.client, req)
}

// Post performs a generic http post on the api
func (c *Connection) Post(endpoint string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, c.hostname+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.makeRequest(endpoint, c.client, req)
}

// Put performs a generic http put on the api
func (c *Connection) Put(endpoint string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, c.hostname+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.makeRequest(endpoint, c.client, req)
}

func (c *Connection) getAndUnmarshal(endpoint string, receive interface{}) error {
	response, err := c.Get(endpoint)
	if err != nil {
		return err
	}

	return json.Unmarshal(response, receive)
}

func (c *Connection) postAndUnmarshal(endpoint string, send interface{}, receive interface{}) error {
	data, err := json.Marshal(send)
	if err != nil {
		return err
	}

	response, err := c.Post(endpoint, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(response, receive)
}

func (c *Connection) putAndUnmarshal(endpoint string, send interface{}, receive interface{}) error {
	data, err := json.Marshal(send)
	if err != nil {
		return err
	}

	response, err := c.Put(endpoint, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(response, receive)
}

func (c *Connection) makeRequest(endpoint string, client http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.authHeader)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, newAPIError(req, res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
