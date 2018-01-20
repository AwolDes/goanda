package goanda

import (
	"log"
	"encoding/json"
	"bytes"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
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