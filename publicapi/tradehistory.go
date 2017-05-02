package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TradeHistory []*Trade

type Trade struct {
	GlobalTradeId int64   `json:"globalTradeID"`
	TradeId       int64   `json:"tradeID"`
	Date          int64   // Unix timestamp
	TypeOrder     string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
}

// Poloniex public API implementation of returnTradeHistory command.
//
// API Doc:
// Returns the past 200 trades for a given market, or up to 50,000 trades between
// a range specified in UNIX timestamps by the "start" and "end" GET parameters.
//
// Call: https://poloniex.com/public?command=returnTradeHistory&currencyPair=BTC_NXT&start=1410158341&end=1410499372
//
// Sample output:
//
//  [
//    {
//      "globalTradeID": 2036467,
//      "tradeID": 21387,
//      "date": "2014-09-12 05:21:26",
//      "type": "buy",
//      "rate": "0.00008943",
//      "amount": "1.27241180",
//      "total": "0.00011379"
//    },
//    {
//      "globalTradeID": 2036466,
//      "tradeID": 21386,
//      "date": "2014-09-12 05:21:25",
//      "type": "buy",
//      "rate": "0.00008943",
//      "amount": "1.27241180",
//      "total": "0.00011379"
//    }, ...
//  ]
func (client *Client) GetTradeHistory(currencyPair string, start, end time.Time) (TradeHistory, error) {

	params := map[string]string{
		"command":      "returnTradeHistory",
		"currencyPair": strings.ToUpper(currencyPair),
		"start":        strconv.Itoa(int(start.Unix())),
		"end":          strconv.Itoa(int(end.Unix())),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	var res = make(TradeHistory, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

// GetPast200TradeHistory returns the past 200 trades for a given market (discarding start and end parameters)
//
// Call: https://poloniex.com/public?command=returnTradeHistory&currencyPair=BTC_NXT
//
// Sample output:
//
//  [
//    {
//      "globalTradeID": 93370042,
//      "tradeID": 1013881,
//      "date": "2017-03-26 02:37:36",
//      "type": "buy",
//      "rate": "0.00001343",
//      "amount": "49.94167793",
//      "total": "0.00067071"
//    },
//    {
//      "globalTradeID": 93369958,
//      "tradeID": 1013880,
//      "date": "2017-03-26 02:37:13",
//      "type": "sell",
//      "rate": "0.00001334",
//      "amount": "33.62895816",
//      "total": "0.00044861"
//    }, ...
//  ]
func (client *Client) GetPast200TradeHistory(currencyPair string) (TradeHistory, error) {

	params := map[string]string{
		"command":      "returnTradeHistory",
		"currencyPair": strings.ToUpper(currencyPair),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := make(TradeHistory, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (t *Trade) UnmarshalJSON(data []byte) error {

	type alias Trade
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("time.Parse: %v", err)
	} else {
		t.Date = int64(timestamp.Unix())
	}

	return nil
}
