package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type OpenOrders []OpenOrder

type OpenOrder struct {
	OrderNumber    int64   `json:"orderNumber,string"`
	Type           string  `json:"type"`
	Rate           float64 `json:"rate,string"`
	StartingAmount float64 `json:"startingAmount,string"`
	Amount         float64 `json:"Amount,string"`
	Total          float64 `json:"Total,string"`
	Date           string  `json:"date"`
	Margin         int     `json:"margin"`
}

type AllOpenOrders map[string]OpenOrders

// Poloniex trading API implementation of returnOpenOrders command.
//
// API Doc:
// Returns your open orders for a given market, specified by the "currencyPair"
// POST parameter, e.g. "BTC_XCP". Set "currencyPair" to "all" to return open
// orders for all markets.
//
// Sample output for single market:
//  [
//    {
//      "orderNumber": "120466",
//      "type": "sell",
//      "rate": "0.025",
//      "startingAmount: 0.025"
//      "amount": "100",
//      "total": "2.5"
//      "date": "2017-03-27 04:21:57"
//      "margin": 0
//     },
//     {
//       "orderNumber": "120467",
//       "type": "sell",
//       "rate": "0.04",
//       "amount": "100",
//       "total": "4"
//     }, ...
//  ]
func (client *TradingClient) GetOpenOrders(currencyPair string) (*OpenOrders, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnOpenOrders")
	postParameters.Add("currencyPair", currencyPair)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := OpenOrders{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}

// GetAllOpenOrders returns the open orders for all markets (currencyPair to "all")
//
// Sample output:
//  {
//    "BTC_ETC": [],
//    "BTC_ETH": [
//      {
//        "orderNumber": "257744844301",
//        "type": "buy",
//        "rate": "0.02",
//        "startingAmount": "0.1",
//        "Amount": "0.1",
//        "Total": "0.002",
//        "date": "2017-03-27 04:25:42",
//        "margin": 0
//      }, ...
//    ], ...
//  }
func (client *TradingClient) GetAllOpenOrders() (AllOpenOrders, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnOpenOrders")
	postParameters.Add("currencyPair", "all")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := make(AllOpenOrders)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}
