package goanda

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

)

type StreamingConnection struct {
	*Connection
	streamURL string
}

func NewStreamingConnection(c *Connection) *StreamingConnection {
	streamURL := "https://stream-fxpractice.oanda.com/v3"
	if strings.Contains(c.hostname, "fxtrade") {
		streamURL = "https://stream-fxtrade.oanda.com/v3"
	}

	return &StreamingConnection{
		Connection: c,
		streamURL:  streamURL,
	}
}

func (c *Connection) NewStreamingConnection() *StreamingConnection {
	return NewStreamingConnection(c)
}

func (sc *StreamingConnection) StreamPrices(instruments []string, callback func(PricingStreamResponse)) error {
	endpoint := fmt.Sprintf("/accounts/%s/pricing/stream", sc.accountID)
	url := sc.streamURL + endpoint + "?instruments=" + strings.Join(instruments, "%2C")

	return sc.stream(url, func(data []byte) error {
		var response PricingStreamResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return err
		}
		if response.Type == "" {
			// This might be an error response
			var errorResp struct {
				ErrorMessage string `json:"errorMessage"`
			}
			if err := json.Unmarshal(data, &errorResp); err == nil && errorResp.ErrorMessage != "" {
				return fmt.Errorf("API error: %s", errorResp.ErrorMessage)
			}
		}
		callback(response)
		return nil
	})
}

func (sc *StreamingConnection) StreamTransactions(callback func(TransactionStreamResponse)) error {
	endpoint := fmt.Sprintf("/accounts/%s/transactions/stream", sc.accountID)
	url := sc.streamURL + endpoint

	return sc.stream(url, func(data []byte) error {
		var response TransactionStreamResponse
		err := json.Unmarshal(data, &response)
		if err != nil {
			return err
		}
		callback(response)
		return nil
	})
}

func (sc *StreamingConnection) stream(url string, handler func([]byte) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", sc.authHeader)
	req.Header.Set("Accept-Datetime-Format", "RFC3339")

	resp, err := sc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	done := make(chan struct{})
	errChan := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			// Handle heartbeats
			if strings.HasPrefix(line, "{\"type\":\"HEARTBEAT\"") {
				var heartbeat HeartbeatResponse
				err := json.Unmarshal([]byte(line), &heartbeat)
				if err == nil {
					fmt.Printf("Received heartbeat at %s\n", heartbeat.Time)
				}
				continue
			}

			err := handler([]byte(line))
			if err != nil {
				errChan <- err
				return
			}
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}

type PricingStreamResponse struct {
	Type       string `json:"type"`
	Time       string `json:"time"`
	Instrument string `json:"instrument,omitempty"`
	Bids       []struct {
		Price     string `json:"price"`
		Liquidity int    `json:"liquidity"`
	} `json:"bids,omitempty"`
	Asks []struct {
		Price     string `json:"price"`
		Liquidity int    `json:"liquidity"`
	} `json:"asks,omitempty"`
	CloseoutBid string `json:"closeoutBid,omitempty"`
	CloseoutAsk string `json:"closeoutAsk,omitempty"`
	Status      string `json:"status,omitempty"`
	Tradeable   bool   `json:"tradeable,omitempty"`
}

type TransactionStreamResponse struct {
	Type          string          `json:"type"`
	Time          string          `json:"time"`
	TransactionID string          `json:"transactionID,omitempty"`
	AccountID     string          `json:"accountID,omitempty"`
	BatchID       string          `json:"batchID,omitempty"`
	RequestID     string          `json:"requestID,omitempty"`
	Transaction   json.RawMessage `json:"transaction,omitempty"`
}

type HeartbeatResponse struct {
	Type string `json:"type"`
	Time string `json:"time"`
}
