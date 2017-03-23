package pushapi

import (
	"fmt"
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
	Last          string
	LowestAsk     string
	HighestBid    string
	PercentChange string
	BaseVolume    string
	QuoteVolume   string
	IsFrozen      float64
	DayHigh       string
	DayLow        string
}

type Ticker <-chan Tick

var ticker chan Tick

func (client *PushClient) SubscribeTicker() (Ticker, error) {

	ticker = make(chan Tick)

	handler := func(args []interface{}, kwargs map[string]interface{}) {

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

// func (t *Tick) UnmarshalJSON(buf []byte) error {

//     tmp := []interface{}{
//         &t.CurrencyPair,
//         &t.Last,
//         &t.LowestAsk,
//         &t.HighestBid,
//         &t.PercentChange,
//         &t.BaseVolume,
//         &t.QuoteVolume,
//         &t.IsFrozen,
//         &t.DayHigh,
//         &t.DayLow,
//     }

//     if err := json.Unmarshal(buf, &tmp); err != nil {
//         return err
//     }

//     if got, want := len(tmp), 10; got != want {
//         return fmt.Errorf("wrong number of fields in Tick: %d != %d", got, want)
//     }

//     return nil
// }
