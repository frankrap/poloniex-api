package main

import (
	"log"
	"poloniex"
	"poloniex/pushapi"
)

var client *pushapi.PushClient

func main() {

	var err error
	client, err = pushapi.NewPushClient()

	if err != nil {
		log.Fatal(err)
	}

	printPushTicker()
}

// Print ticker periodically
func printPushTicker() {

	ticker, err := client.SubscribeTicker()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg := <-ticker
			poloniex.PrettyPrintJson(msg)
		}

	}()

	select {}
}
