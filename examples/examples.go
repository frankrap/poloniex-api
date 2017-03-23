package main

import (
	"api/poloniex"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {

	// Print ticker periodically
	/*
		client, err := poloniex.NewPushClient()

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
				fmt.Println(poloniex.Tick(msg).Last)
			}

		}()

		select {}
	*/

	client := poloniex.NewPublicClient()

	// Print All order books with depth 2
	/*

		res, err := client.GetAllOrderBooks(2)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)
	*/

	// Print BTC_STEEM order book with depth 200
	/*
		ob, _ := client.GetOrderBook("BTC_STEEM", 200)
		fmt.Println(len(ob.Asks))
	*/
}
