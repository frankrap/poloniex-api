package main

import (
	"fmt"
	"log"
	"time"
	"trading/poloniex"
	"trading/poloniex/pushapi"
)

var client *pushapi.PushClient

func main() {

	var err error
	client, err = pushapi.NewPushClient()

	if err != nil {
		log.Fatal(err)
	}

	// go printTicker()
	// go printTrollbox()
	go printMarketUpdates()
	select {}
}

// Print ticker periodically
func printTicker() {

	done := make(chan struct{})
	ticker, err := client.SubscribeTicker()
	if err != nil {
		log.Fatal(err)
	}

	loop := func() {

		for {
			select {
			case msg := <-ticker:
				poloniex.PrettyPrintJson(msg)
			case <-done:
				//return
			}
		}
	}

	go loop()

	time.Sleep(3 * time.Second)
	client.UnsubscribeTicker()
	done <- struct{}{}

	time.Sleep(3 * time.Second)
	client.SubscribeTicker()

	time.Sleep(3 * time.Second)
	client.UnsubscribeTicker()
	done <- struct{}{}
}

// Print trollbox periodically
func printTrollbox() {

	done := make(chan struct{})

	trollbox, err := client.SubscribeTrollbox()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case msg := <-trollbox:
				fmt.Printf("%d | %s: %s\n", msg.Reputation, msg.Username, msg.Message)
			case <-done:
				return
			}
		}

	}()

	time.Sleep(15 * time.Second)
	client.UnsubscribeTrollbox()
	done <- struct{}{}
}

func printMarketUpdates() {

	done := make(chan struct{})
	marketUpdate, err := client.SubscribeMarket("BTC_ETH")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case msg := <-marketUpdate:
				poloniex.PrettyPrintJson(msg)
			case <-done:
				return
			}
		}
	}()
	time.Sleep(3 * time.Second)
	client.UnsubscribeMarket("BTC_ETH")
	client.SubscribeMarket("BTC_ETH")

	time.Sleep(2 * time.Second)
	client.UnsubscribeMarket("BTC_ETH")
	done <- struct{}{}
}
