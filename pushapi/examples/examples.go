package main

import (
	"fmt"
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

	go printTicker()
	go printTrollbox()
	select {}
}

// Print ticker periodically
func printTicker() {

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
		time.Sleep(1 * time.Second)
		client.UnsubscribeTicker()
	}()
}

// Print trollbox periodically
func printTrollbox() {

	trollbox, err := client.SubscribeTrollbox()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msg, ok := <-trollbox
			if !ok {
				break
			}
			fmt.Printf("%d | %s: %s\n", msg.Reputation, msg.Username, msg.Message)
		}

	}()

	go func() {
		time.Sleep(5 * time.Second)
		client.UnsubscribeTrollbox()
	}()
}
