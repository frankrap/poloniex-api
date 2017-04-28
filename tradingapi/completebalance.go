package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type CompleteBalances map[string]*CompleteBalance

type CompleteBalance struct {
	Available float64 `json:"available,string"`
	OnOrders  float64 `json:"onOrders,string"`
	BtcValue  float64 `json:"btcValue,string"`
}

// Poloniex trading API implementation of returnCompleteBalances command.
//
// API Doc:
// Returns all of your balances, including available balance, balance on orders,
// and the estimated BTC value of your balance. By default, this call is limited
// to your exchange account; set the "account" POST parameter to "all" to include
// your margin and lending accounts.
//
// Sample output:
//
//  {
//    "LTC": {
//      "available": "5.015",
//      "onOrders": "1.0025",
//      "btcValue": "0.078"
//    }, ...
//  }
func (client *Client) GetCompleteBalances() (CompleteBalances, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnCompleteBalances")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(CompleteBalances)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}
