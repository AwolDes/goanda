# goanda
An Golang wrapper for the [OANDA](http://developer.oanda.com/) v20 API.

## Requirements
- Go v1.9+

_Note: This package was created by a third party, and was not created by anyone affiliated with OANDA_

## Usage
To use this package run `go get github.com/awoldes/goanda` then import it into your program and set it up with the following snippet:
```
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/awoldes/goanda"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	oanda := goanda.NewConnection(key, false)
	history := oanda.Request("v3/instruments/EUR_USD/candles")
	fmt.Println(history)
}

```
