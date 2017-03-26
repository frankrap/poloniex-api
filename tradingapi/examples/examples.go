package main

import (
	"log"
	"poloniex"
	"poloniex/tradingapi"
)

var client *tradingapi.TradingClient

func main() {

	var err error
	client, err = tradingapi.NewTradingClient(poloniex.API_KEY, poloniex.API_SECRET)

	if err != nil {
		log.Fatal(err)
	}

	// printBalances()
	printCompleteBalances()
}

// Print balances
func printBalances() {

	res, err := client.GetBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print complete balances
func printCompleteBalances() {

	res, err := client.GetCompleteBalances()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}
