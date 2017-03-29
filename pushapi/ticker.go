package pushapi

import (
	"fmt"
	"strconv"
)

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

type Ticker chan *Tick

var ticker Ticker

// Poloniex push API implementation of ticker topic.
//
// API Doc:
// In order to receive ticker updates, subscribe to "ticker".
//
// Updates will be in the following format:
//
// [currencyPair, last, lowestAsk, highestBid, percentChange,
//  baseVolume, quoteVolume, isFrozen, 24hrHigh, 24hrLow]
//
// Example:
//
// ['BTC_BBR','0.00069501','0.00074346','0.00069501', '-0.00742634',
//  '8.63286802','11983.47150109',0,'0.00107920','0.00045422']
func (client *PushClient) SubscribeTicker() (Ticker, error) {

	ticker = make(Ticker)

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		if tick, err := convertArgsToTick(args); err != nil {
			fmt.Printf("convertArgstoTick: %v\n", err)
		} else {
			ticker <- tick
		}
	}

	if err := client.wampClient.Subscribe(TICKER, nil, handler); err != nil {
		return nil, fmt.Errorf("turnpike.Client.Subscribe: %v", err)
	}

	return ticker, nil
}

func (client *PushClient) UnsubscribeTicker() error {

	if err := client.wampClient.Unsubscribe(TICKER); err != nil {
		return fmt.Errorf("turnpike.Client.Unsuscribe: %v", err)
	}
	close(ticker)

	return nil
}

func convertArgsToTick(args []interface{}) (*Tick, error) {

	convertArg := func(arg interface{}) (float64, error) {

		if v, ok := arg.(string); ok {

			val, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0, fmt.Errorf("strconv.ParseFloat: %v", err)
			}
			return val, nil

		} else {
			return 0, fmt.Errorf("type assertion failed: %v", arg)
		}
	}

	var tick = Tick{}
	var err error

	if v, ok := args[0].(string); ok {
		tick.CurrencyPair = v
	} else {
		return nil, fmt.Errorf("'CurrencyPair' type assertion failed")
	}

	if tick.Last, err = convertArg(args[1]); err != nil {
		return nil, fmt.Errorf("convertArg 'Last': %v", err)
	} else if tick.LowestAsk, err = convertArg(args[2]); err != nil {
		return nil, fmt.Errorf("convertArg 'LowestAsk': %v", err)
	} else if tick.HighestBid, err = convertArg(args[3]); err != nil {
		return nil, fmt.Errorf("convertArg 'HighestBid': %v", err)
	} else if tick.PercentChange, err = convertArg(args[4]); err != nil {
		return nil, fmt.Errorf("convertArg 'PercentChange': %v", err)
	} else if tick.BaseVolume, err = convertArg(args[5]); err != nil {
		return nil, fmt.Errorf("convertArg 'BaseVolume': %v", err)
	} else if tick.QuoteVolume, err = convertArg(args[6]); err != nil {
		return nil, fmt.Errorf("convertArg 'QuoteVolume': %v", err)
	}

	if v, ok := args[7].(float64); ok {
		if v == 0 {
			tick.IsFrozen = false
		} else {
			tick.IsFrozen = true
		}
	} else {
		return nil, fmt.Errorf("'IsFrozen' type assertion failed")
	}

	if tick.High24hr, err = convertArg(args[8]); err != nil {
		return nil, fmt.Errorf("convertArg 'High24hr': %v", err)
	} else if tick.Low24hr, err = convertArg(args[9]); err != nil {
		return nil, fmt.Errorf("convertArg 'Low24hr': %v", err)
	}

	return &tick, nil
}
