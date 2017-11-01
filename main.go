package main

import (
	"fmt"
	"go_oanda"
)

func main() {
	key := "OANDA_API_KEY"
	oanda := go_oanda.NewConnection(key, false)
	history := oanda.Request("v3/instruments/EUR_USD/candles")
	fmt.Println(history)
}
