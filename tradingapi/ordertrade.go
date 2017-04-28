package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type TradesFromOrder []*TradeFromOrder

type TradeFromOrder struct {
	GlobalTradeId int64   `json:"globalTradeID"`
	TradeId       int64   `json:"tradeID"`
	CurrencyPair  string  `json:"currencyPair"`
	TypeOrder     string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
	Fee           float64 `json:"fee,string"`
	Date          int64   // Unix timestamp
}

// Poloniex trading API implementation of returnOrderTrades command.
//
// API Doc:
// Returns all trades involving a given order, specified by the "orderNumber" POST parameter.
// If no trades for the order have occurred or you specify an order that does not belong to you,
// you will receive an error.
//
// Sample output:
//
//  [
//    {
//      "globalTradeID": 89366140,
//      "tradeID": 652357,
//      "currencyPair": "BTC_STEEM",
//      "type": "buy",
//      "rate": "0.00021999",
//      "amount": "53.30947121",
//      "total": "0.01172755",
//      "fee": "0.0025",
//      "date": "2017-03-18 06:28:20"
//    },
//    {
//      "globalTradeID": 89366139,
//      "tradeID": 652356,
//      "currencyPair": "BTC_STEEM",
//      "type": "buy",
//      "rate": "0.00021998",
//      "amount": "1.02657424",
//      "total": "0.00022582",
//      "fee": "0.0025"
//      "date": "2017-03-18 06:28:20"
//    }
//  ]
func (client *Client) GetTradesFromOrder(orderNumber int64) (TradesFromOrder, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnOrderTrades")
	postParameters.Add("orderNumber", strconv.Itoa(int(orderNumber)))

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(TradesFromOrder, 0)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (t *TradeFromOrder) UnmarshalJSON(data []byte) error {

	type alias TradeFromOrder
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
