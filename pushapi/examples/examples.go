package main

import (
	"log"
	"poloniex"
	"poloniex/pushapi"
	"time"
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
			msg, ok := <-ticker
			if !ok {
				break
			}
			poloniex.PrettyPrintJson(msg)
		}

	}()

	go func() {
		time.Sleep(2 * time.Second)
		client.UnsubscribeTicker()
	}()

	select {}
}
