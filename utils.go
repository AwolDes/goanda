package goanda

import (
	"log"
	"encoding/json"
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