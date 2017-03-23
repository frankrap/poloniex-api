package pushapi

import (
	"fmt"
	"log"
	"strconv"
)

// MESSAGE FORMAT:
// [currencyPair, last, lowestAsk, highestBid, percentChange,
//  baseVolume, quoteVolume, isFrozen, 24hrHigh, 24hrLow]
//
// Example:
// ['BTC_BBR','0.00069501','0.00074346','0.00069501', '-0.00742634',
//  '8.63286802','11983.47150109',0,'0.00107920','0.00045422']

const (
	TICKER = "ticker"
)

type Tick struct {
	CurrencyPair  string
	Last          float64
	LowestAsk     float64
	HighestBid    float64
	PercentChange float64
	BaseVolume    float64
	QuoteVolume   float64
	IsFrozen      bool
	High24hr      float64
	Low24hr       float64
}

type Ticker <-chan Tick

var ticker chan Tick

func (client *PushClient) SubscribeTicker() (Ticker, error) {

	ticker = make(chan Tick)

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		var tick = Tick{}

		tick.CurrencyPair = args[0].(string)

		val, err := strconv.ParseFloat(args[1].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.Last = val

		val, err = strconv.ParseFloat(args[2].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.LowestAsk = val

		val, err = strconv.ParseFloat(args[3].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.HighestBid = val

		val, err = strconv.ParseFloat(args[4].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.PercentChange = val

		val, err = strconv.ParseFloat(args[5].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.BaseVolume = val

		val, err = strconv.ParseFloat(args[6].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.QuoteVolume = val

		f := args[7].(float64)
		if f == 0 {
			tick.IsFrozen = false
		} else {
			tick.IsFrozen = true
		}

		val, err = strconv.ParseFloat(args[8].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.High24hr = val

		val, err = strconv.ParseFloat(args[9].(string), 64)
		if err != nil {
			log.Fatal(err)
		}
		tick.Low24hr = val

		ticker <- tick
	}

	if err := client.wampClient.Subscribe("ticker", nil, handler); err != nil {
		return nil, fmt.Errorf("subscribe %s: %v", TICKER, err)
	}

	return ticker, nil
}

func (client *PushClient) UnsubscribeTicker() error {

	if err := client.wampClient.Unsubscribe(TICKER); err != nil {
		return fmt.Errorf("unsuscribe %s: %v", TICKER, err)
	}

	close(ticker)
	return nil
}
