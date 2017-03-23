package pushapi

import (
	"fmt"
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

	convertArg := func(arg interface{}) (float64, error) {

		if v, ok := arg.(string); ok {

			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, fmt.Errorf(" : %v", arg)
			}
			return val, nil

		} else {
			return 0, fmt.Errorf("type assertion failed: %v", arg)
		}
	}

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		var tick = Tick{}

		if v, ok := args[0].(string); ok {
			tick.CurrencyPair = v
		} else {
			return
		}

		var err error

		if tick.Last, err = convertArg(args[1]); err != nil {
			return
		}
		if tick.LowestAsk, err = convertArg(args[2]); err != nil {
			return
		}
		if tick.HighestBid, err = convertArg(args[3]); err != nil {
			return
		}
		if tick.PercentChange, err = convertArg(args[4]); err != nil {
			return
		}
		if tick.BaseVolume, err = convertArg(args[5]); err != nil {
			return
		}
		if tick.QuoteVolume, err = convertArg(args[6]); err != nil {
			return
		}

		if v, ok := args[7].(float64); ok {
			if v == 0 {
				tick.IsFrozen = false
			} else {
				tick.IsFrozen = true
			}
		} else {
			return
		}

		if tick.High24hr, err = convertArg(args[8]); err != nil {
			return
		}
		if tick.Low24hr, err = convertArg(args[9]); err != nil {
			return
		}

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
