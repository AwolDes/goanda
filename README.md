# goanda
A Golang wrapper for the [OANDA v20 API](http://developer.oanda.com/rest-live-v20/introduction/). Currently OANDA has wrappers for Python, Javascript and Java. Goanda exists to extend upon those languages because of the increasing popularity of Go and for a side prject I'm working on.

## Requirements
- Go v1.9+

_Note: This package was created by a third party, and was not created by anyone affiliated with OANDA_

## Usage
To use this package run `go get github.com/awoldes/goanda` then import it into your program and set it up following the snippets below.

### Price History
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
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)
	history := oanda.GetCandles("EUR_USD")
	fmt.Println(history)
}

```

Look at the `/examples` directory for more!

## TODO
### **API** (in order of priority)
- [x] Instrument endpoints (to get prices and the order book)
- [x] Order endpoints (to create, get or update orders for an account)
- [ ] Trade enpoints (to get information on current trades) 
- [ ] Position endpoints (to get information on current positions)
- [ ] Account endpoints (to get information on the account)
- [ ] Pricing endpoints (to get pricing of instruments or stream)
- [ ] Transaction endpoints (to get information on account transactions)

### **Docs**
- [ ] Write docs on how to use `goanda`
- [ ] Write example programs for `goanda`
