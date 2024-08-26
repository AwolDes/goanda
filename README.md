# goanda
A Golang wrapper for the [OANDA v20 API](http://developer.oanda.com/rest-live-v20/introduction/). Currently OANDA has wrappers for Python, Javascript and Java. Goanda exists to extend upon those languages because of the increasing popularity of Go.

## Features
Goanda can perform the following actions on your OANDA brokerage accounts:

- Get candlesticks of all instruments
- Create and update orders with comprehensive order properties
- Get data on current and past trades on OANDA
- Close/scale out of trades you have open
- Close positions (not just trades)
- Get data on your account
  - NAV
  - Current % of used margin
  - Balance
  - And more!
- Get data on all your transactions
- Get all pricing data (bid/ask spread) on specific instruments
- Stream real-time data for prices, transactions, account changes, and candles

## Requirements
- Go v1.9+

_Note: This package was created by a third party, and was not created by anyone affiliated with OANDA_

## Usage
To use this package run `go get github.com/awoldes/goanda` then import it into your program and set it up following the snippets below.

### Basic Example
I suggest creating a `.env` file for your project to keep your secrets safe! Make sure you add a .gitignore file.

`~/project/.env`

```
OANDA_API_KEY=
OANDA_ACCOUNT=
```

`~/project/main.go`

```
package main

import (
	"log"
	"os"

	"github.com/awoldes/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &goanda.ConnectionConfig{
		UserAgent: "goanda",
		Timeout: 10 * time.Second,
		Live: false,
	}

	granularity := goanda.GranularityFiveSeconds

	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")

	oanda, err := goanda.NewConnection(accountID, key, config)
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}

	history, err := oanda.GetCandles("EUR_USD", 10, granularity)
	if err != nil {
		log.Fatalf("Error getting candles: %v", err)
	}

	spew.Dump(history)
}

```

Look at the [`/examples`](https://github.com/AwolDes/goanda/tree/master/examples) directory for more!

## Contributing
For now if you'd like to contribute create an Issue and/or submit a PR!

## Testing
The project now includes comprehensive test coverage. To run the tests, use the following command:

```
go test -v ./...
```

## TODO
### **API** (in order of priority)
- [x] Instrument endpoints (to get prices and the order book)
- [x] Order endpoints (to create, get or update orders for an account)
- [x] Trade endpoints (to get information on current trades) 
- [x] Position endpoints (to get information on current positions)
- [x] Account endpoints (to get information on the account)
- [x] Transaction endpoints (to get information on account transactions)
- [x] Pricing endpoints (to get pricing of instruments)
- [x] Streaming endpoints for Pricing & Transactions

### **Docs**
- [x] Write docs on how to use `goanda`
- [x] Write example programs for `goanda`
- [x] Write tests for `goanda`


## Supporting Projects
Thank you to the following projects, they really helped me while I was developing this API

- [JSON to Go from @mholt](https://mholt.github.io/json-to-go/) - this tool helped massively with all the structs that needed to be created.
- [Spew from @davecgh](https://github.com/davecgh/go-spew) - this tool was awesome to debug with.

## License

This project was created under the [MIT license](https://choosealicense.com/licenses/mit/)
