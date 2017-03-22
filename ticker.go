package poloniex

import (
	"fmt"

	turnpike "gopkg.in/jcelliott/turnpike.v2"
)

// MESSAGE FORMAT:
// [currencyPair, last, lowestAsk, highestBid, percentChange,
//  baseVolume, quoteVolume, isFrozen, 24hrHigh, 24hrLow]
//
// Example:
// ['BTC_BBR','0.00069501','0.00074346','0.00069501', '-0.00742634',
//  '8.63286802','11983.47150109',0,'0.00107920','0.00045422']

type Tick struct {
	currencyPair  string
	last          string
	lowestAsk     string
	highestBid    string
	percentChange string
	baseVolume    string
	quoteVolume   string
	isFrozen      float64
	dayHigh       string
	dayLow        string
}

type Ticker <-chan Tick

var ticker chan Tick

func SubscribeTicker(client *turnpike.Client) (Ticker, error) {

	ticker = make(chan Tick)

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		// fmt.Println(args)
		ticker <- Tick{
			args[0].(string),
			args[1].(string),
			args[2].(string),
			args[3].(string),
			args[4].(string),
			args[5].(string),
			args[6].(string),
			args[7].(float64),
			args[8].(string),
			args[9].(string),
		}
	}

	if err := client.Subscribe("ticker", nil, handler); err != nil {
		return nil, fmt.Errorf("subscribe %s: %v", TICKER, err)
	}

	return ticker, nil
}

func UnsubscribeTicker(client *turnpike.Client) error {

	if err := client.Unsubscribe(TICKER); err != nil {
		return fmt.Errorf("unsuscribe %s: %v", TICKER, err)
	}

	close(ticker)
	return nil
}
