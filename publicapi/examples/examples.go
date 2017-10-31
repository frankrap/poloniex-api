package main

import (
	"log"
	"time"

	poloniex "github.com/joemocquant/poloniex-api"
	publicapi "github.com/joemocquant/poloniex-api/publicapi"
)

var client *publicapi.Client

// go run example.go
func main() {

	client = publicapi.NewClient()

	printPublicTickers()

	// printPublicDayVolumes()

	// printPublicOrderBook()

	// printPublicOrderBooks()

	// printTradeHistory()

	// printPast200TradeHistory()

	// printChartData()

	// printCurrencies()

	// printLoanOrders()
}

func printPublicTickers() {

	res, err := client.GetTickers()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

func printPublicDayVolumes() {

	res, err := client.GetDayVolumes()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print BTC_STEEM order book with depth 200
func printPublicOrderBook() {

	res, err := client.GetOrderBook("BTC_STEEM", 200)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print All order books with depth 2
func printPublicOrderBooks() {

	res, err := client.GetOrderBooks(2)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print BTC_STEEM trade the last 10 minutes
func printTradeHistory() {

	end := time.Now()
	start := end.Add(-10 * time.Minute)
	res, err := client.GetTradeHistory("BTC_STEEM", start, end)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print past 200 BTC_STEEM trades
func printPast200TradeHistory() {

	res, err := client.GetPast200TradeHistory("BTC_STEEM")

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print BTC_STEEM  30min candlesticks the last 10 hours
func printChartData() {

	end := time.Now()
	start := end.Add(-10 * time.Hour)
	res, err := client.GetChartData("BTC_STEEM", start, end, 1800)

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

func printCurrencies() {

	res, err := client.GetCurrencies()

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}

// Print loan orders for BTC
func printLoanOrders() {

	res, err := client.GetLoanOrders("BTC")

	if err != nil {
		log.Fatal(err)
	}

	poloniex.PrettyPrintJson(res)
}
