package main

import (
	"api/poloniex/publicapi"
	"api/poloniex/pushapi"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {

	// printPushTicker()
	// printPublicAllTick()
	// printPublicAllDayVolumes()
	// printPublicAllOrderBook()
	// printPublicOrderBook()
	// printPast200TradeHistory()
	printTradeHistory()
}

func prettyPrintJson(msg interface{}) {
	jsonstr, err := json.MarshalIndent(msg, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

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

func printPublicAllTick() {

	client := publicapi.NewPublicClient()

	res, err := client.GetTicker()

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

func printPublicAllDayVolumes() {

	client := publicapi.NewPublicClient()

	res, err := client.GetAllDayVolumes()

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

func printPast200TradeHistory() {

	client := publicapi.NewPublicClient()

	res, err := client.GetPast200TradeHistory("BTC_STEEM")

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}

// Print BTC_STEEM trade the last 10 minutes
func printTradeHistory() {

	client := publicapi.NewPublicClient()

	end := time.Now()
	start := end.Add(-10 * time.Minute)
	res, err := client.GetTradeHistory("BTC_STEEM", start, end)

	if err != nil {
		log.Fatal(err)
	}

	prettyPrintJson(res)
}
