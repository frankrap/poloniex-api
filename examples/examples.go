package main

import (
	"api/poloniex/publicapi"
	"api/poloniex/pushapi"
	"encoding/json"
	"fmt"
	"log"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {

	// printPushTicker()
	printPublicTick()
	// printPublicAllOrderBook()
	// printPublicOrderBook()
}

func prettyPrintJson(msg interface{}) {
	jsonstr, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Printf("%s\n", string(jsonstr))
}

//
// PUSH API
//

// Print ticker periodically
func printPushTicker() {
	client, err := pushapi.NewPushClient()

	if err != nil {
		log.Fatal(err)
	}

	ticker, err := client.SubscribeTicker()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-ticker
			prettyPrintJson(msg)
		}

	}()

	select {}
}

//
// PUBLIC API
//

func printPublicTick() {
	client := publicapi.NewPublicClient()

	res, err := client.GetTicker()

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

// Print All order books with depth 2
func printPublicAllOrderBook() {

	client := publicapi.NewPublicClient()

	res, err := client.GetAllOrderBooks(2)

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

// Print BTC_STEEM order book with depth 200
func printPublicOrderBook() {

	client := publicapi.NewPublicClient()

	res, err := client.GetOrderBook("BTC_STEEM", 200)

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}
