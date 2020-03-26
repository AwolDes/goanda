package goanda

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func checkApiErr(body []byte, route string) error {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		return errors.New("\nOANDA API Error: " + bodyString + "\nOn route: " + route)
	}
	return nil
}

func unmarshalJson(body []byte, data interface{}) error {
	return json.Unmarshal(body, &data)
}

func createUrl(host string, endpoint string) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)

	url := buffer.String()
	return url
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, getErr := client.Do(req)
	if getErr !=nil {
		log.Printf("Error occurred in client.Do(): %v\n", getErr)
		return []byte{}, getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Printf("Error occurred in ioutil.ReadAll(): %v\n", readErr)
		return []byte{}, readErr
	}
	apiErr := checkApiErr(body, endpoint)
	if apiErr != nil {
		return []byte{}, apiErr
	}
	return body, nil
}
