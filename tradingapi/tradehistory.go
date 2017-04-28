package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type TradeHistory []*Trade

type Trade struct {
	GlobalTradeId int64   `json:"globalTradeID"`
	TradeId       int64   `json:"tradeID,string"`
	Date          int64   // Unix timestamp
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
	Fee           float64 `json:"fee,string"`
	OrderNumber   int64   `json:"orderNumber,string"`
	TypeOrder     string  `json:"type"`
	Category      string  `json:"category"`
}

type AllTradeHistory map[string]TradeHistory

// Poloniex trading API implementation of returnTradeHistory command.
//
// API Doc:
// Returns your trade history for a given market, specified by the "currencyPair" POST parameter.
// You may specify "all" as the currencyPair to receive your trade history for all markets. You
// may optionally specify a range via "start" and/or "end" POST parameters, given in UNIX
// timestamp format; if you do not specify a range, it will be limited to one day.
//
// Sample output:
//
//  [
//    {
//      "globalTradeID": 25129732,
//      "tradeID": "6325758",
//      "date": "2016-04-05 08:08:40",
//      "rate": "0.02565498",
//      "amount": "0.10000000",
//      "total": "0.00256549",
//      "fee": "0.00200000",
//      "orderNumber": "34225313575",
//      "type": "sell",
//      "category": "exchange"
//    },
//    {
//      "globalTradeID": 25129628,
//      "tradeID": "6325741",
//      "date": "2016-04-05 08:07:55",
//      "rate": "0.02565499",
//      "amount": "0.10000000",
//      "total": "0.00256549",
//      "fee": "0.00200000",
//      "orderNumber": "34225195693",
//      "type": "buy",
//      "category": "exchange"
//    }, ...
//  ]
func (client *Client) GetTradeHistory(currencyPair string, start, end time.Time) (TradeHistory, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnTradeHistory")
	postParameters.Add("currencyPair", currencyPair)
	postParameters.Add("start", strconv.Itoa(int(start.Unix())))
	postParameters.Add("end", strconv.Itoa(int(end.Unix())))

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(TradeHistory, 0)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

// GetAllOpenOrders returns the open orders for all markets (currencyPair to "all")
//
// Sample output:
//
//  { "BTC_ETH": [
//      {
//        "globalTradeID": 25129732,
//        "tradeID": "6325758",
//        "date": "2016-04-05 08:08:40",
//        "rate": "0.02565498",
//        "amount": "0.10000000",
//        "total": "0.00256549",
//        "fee": "0.00200000",
//        "orderNumber": "34225313575",
//        "type": "sell",
//        "category": "exchange"
//      },
//      {
//        "globalTradeID": 25129628,
//        "tradeID": "6325741",
//        "date": "2016-04-05 08:07:55",
//        "rate": "0.02565499",
//        "amount": "0.10000000",
//        "total": "0.00256549",
//        "fee": "0.00200000",
//        "orderNumber": "34225195693",
//        "type": "buy",
//        "category": "exchange"
//      }, ...
//    ], ...
//  }
func (client *Client) GetAllTradeHistory(start, end time.Time) (AllTradeHistory, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnTradeHistory")
	postParameters.Add("currencyPair", "all")
	postParameters.Add("start", strconv.Itoa(int(start.Unix())))
	postParameters.Add("end", strconv.Itoa(int(end.Unix())))

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(AllTradeHistory, 0)

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
