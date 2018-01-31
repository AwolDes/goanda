package goanda

import (
	"bytes"
	"encoding/json"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func unmarshalJson(body []byte, data interface{}) {
	// TODO: Better way to handle error responses from oanda
	// I.e: {"errorMessage":"Invalid value specified for 'accountID'"}
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
