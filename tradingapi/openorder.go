package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type OpenOrders []*OpenOrder

type OpenOrder struct {
	OrderNumber    int64   `json:"orderNumber,string"`
	Type           string  `json:"type"`
	Rate           float64 `json:"rate,string"`
	StartingAmount float64 `json:"startingAmount,string"`
	Amount         float64 `json:"Amount,string"`
	Total          float64 `json:"Total,string"`
	Date           int64   // Unix timestamp
	Margin         int     `json:"margin"`
}

type AllOpenOrders map[string]*OpenOrders

// Poloniex trading API implementation of returnOpenOrders command.
//
// API Doc:
// Returns your open orders for a given market, specified by the "currencyPair"
// POST parameter, e.g. "BTC_XCP". Set "currencyPair" to "all" to return open
// orders for all markets.
//
// Sample output for single market:
//
//  [
//    {
//      "orderNumber": "258029798062",
//      "type": "buy",
//      "rate": "0.0048671",
//      "startingAmount": "0.1",
//      "Amount": "0.1",
//      "Total": "0.00048671",
//      "date": "2017-03-27 16:46:16",
//      "margin": 0
//    },
//    {
//      "orderNumber": "258029833027",
//      "type": "buy",
//      "rate": "0.0048671",
//      "startingAmount": "0.1",
//      "Amount": "0.1",
//      "Total": "0.00048671",
//      "date": "2017-03-27 16:46:21",
//      "margin": 0
//    }, ...
//  ]
func (client *Client) GetOpenOrders(currencyPair string) (*OpenOrders, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnOpenOrders")
	postParameters.Add("currencyPair", currencyPair)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := OpenOrders{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

// GetAllOpenOrders returns the open orders for all markets (currencyPair to "all")
//
// Sample output:
//
//  {
//    "BTC_ETC": [],
//    "BTC_ETH": [
//      {
//        "orderNumber": "258029798062",
//        "type": "buy",
//        "rate": "0.0048671",
//        "startingAmount": "0.1",
//        "Amount": "0.1",
//        "Total": "0.00048671",
//        "date": "2017-03-27 16:46:16",
//        "margin": 0
//      },
//      {
//        "orderNumber": "258029833027",
//        "type": "buy",
//        "rate": "0.0048671",
//        "startingAmount": "0.1",
//        "Amount": "0.1",
//        "Total": "0.00048671",
//        "date": "2017-03-27 16:46:21",
//        "margin": 0
//      }, ...
//    ], ...
//  }
func (client *Client) GetAllOpenOrders() (AllOpenOrders, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnOpenOrders")
	postParameters.Add("currencyPair", "all")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := make(AllOpenOrders)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (o *OpenOrder) UnmarshalJSON(data []byte) error {

	type alias OpenOrder
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("time.Parse: %v", err)
	} else {
		o.Date = int64(timestamp.Unix())
	}

	return nil
}
