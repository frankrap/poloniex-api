package pushapi

import (
	"errors"
	"fmt"
	"sync"
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

var (
	ticker = make(Ticker)

	tickerMu           sync.RWMutex
	tickerUnsubscribed = make(chan struct{})
)

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
func (client *Client) SubscribeTicker() (Ticker, error) {

	handler := func(args []interface{}, kwargs map[string]interface{}) {

		client.updateTopicTimestamp(TICKER)

		tick, err := convertArgsToTick(args)
		if err != nil {
			logger.WithField("error", err).Error("convertArgstoTick")
			return
		}

		tickerMu.RLock()
		select {
		case ticker <- tick:
		case <-tickerUnsubscribed:
		}
		tickerMu.RUnlock()
	}

	subscribe := func() error {

		client.wampClientMu.RLock()
		defer client.wampClientMu.RUnlock()

		if err := client.wampClient.Subscribe(TICKER, nil, handler); err != nil {
			return fmt.Errorf("turnpike.Client.Subscribe: %v", err)
		}
		logger.Infof("Subscribed to: %s", TICKER)

		return nil
	}

	if err := subscribe(); err != nil {
		return nil, err
	}

	tickerMu.Lock()
	select {
	case <-tickerUnsubscribed:
		tickerUnsubscribed = make(chan struct{})
	default:
	}
	tickerMu.Unlock()

	client.addSubscription(TICKER, subscribe)

	return ticker, nil
}

func (client *Client) UnsubscribeTicker() error {

	client.wampClientMu.RLock()

	if err := client.wampClient.Unsubscribe(TICKER); err != nil {
		return fmt.Errorf("turnpike.Client.Unsuscribe: %v", err)
	}

	client.wampClientMu.RUnlock()

	client.removeSubscription(TICKER)

	ticker <- nil

	tickerMu.RLock()
	close(tickerUnsubscribed)
	tickerMu.RUnlock()

	return nil
}

func convertArgsToTick(args []interface{}) (*Tick, error) {

	var tick = Tick{}
	var err error

	if v, ok := args[0].(string); ok {
		tick.CurrencyPair = v
	} else {
		return nil, fmt.Errorf("'CurrencyPair' type assertion failed")
	}

	if tick.Last, err = convertStringToFloat(args[1]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'Last': %v", err)
	} else if tick.LowestAsk, err = convertStringToFloat(args[2]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'LowestAsk': %v", err)
	} else if tick.HighestBid, err = convertStringToFloat(args[3]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'HighestBid': %v", err)
	} else if tick.PercentChange, err = convertStringToFloat(args[4]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'PercentChange': %v", err)
	} else if tick.BaseVolume, err = convertStringToFloat(args[5]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'BaseVolume': %v", err)
	} else if tick.QuoteVolume, err = convertStringToFloat(args[6]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'QuoteVolume': %v", err)
	}

	if v, ok := args[7].(float64); ok {
		if v == 0 {
			tick.IsFrozen = false
		} else {
			tick.IsFrozen = true
		}
	} else {
		return nil, errors.New("'IsFrozen' type assertion failed")
	}

	if tick.High24hr, err = convertStringToFloat(args[8]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'High24hr': %v", err)
	} else if tick.Low24hr, err = convertStringToFloat(args[9]); err != nil {
		return nil, fmt.Errorf("convertStringToFloat 'Low24hr': %v", err)
	}

	return &tick, nil
}
