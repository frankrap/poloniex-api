package main

import (
	"fmt"
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
	// printCompleteBalances()
	printDepositAddresses()
	GenerateNewAddress()
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

// Print deposit addresses
func printDepositAddresses() {

	res, err := client.GetDepositAddresses()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Generate new address for currency BTC
func GenerateNewAddress() {

	currency := "BTC"
	addr, err := client.GenerateNewAddress(currency)

	if err != nil {
		log.Fatal(err)
	}

	toPrint := fmt.Sprintf("New address generated (%s): %s", currency, addr)
	poloniex.PrettyPrintJson(toPrint)
}
