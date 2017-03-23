package main

import (
	"api/poloniex/publicapi"
	"api/poloniex/pushapi"
	"fmt"
	"log"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	printTicker()
}

// Print ticker periodically
func printTicker() {
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
			fmt.Printf("%#v\n", pushapi.Tick(msg))
		}

	}()

	select {}
}

// Print All order books with depth 2
func printAllOrderBook() {

	client := publicapi.NewPublicClient()

	res, err := client.GetAllOrderBooks(2)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

// Print BTC_STEEM order book with depth 200
func printOrderBook() {

	client := publicapi.NewPublicClient()
	ob, _ := client.GetOrderBook("BTC_STEEM", 200)
	fmt.Println(len(ob.Asks))

}
