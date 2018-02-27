package goanda

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkApiErr(body []byte, route string) {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		log.SetFlags(log.LstdFlags | log.Llongfile)
		log.Fatal("\nOANDA API Error: " + bodyString + "\nOn route: " + route)
	}
}

func unmarshalJson(body []byte, data interface{}) {
	jsonErr := json.Unmarshal(body, &data)
	checkErr(jsonErr)
}

func createUrl(host string, endpoint string) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)

	url := buffer.String()
	return url
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) []byte {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, getErr := client.Do(req)
	checkErr(getErr)
	body, readErr := ioutil.ReadAll(res.Body)
	checkErr(readErr)
	checkApiErr(body, endpoint)
	return body
}
