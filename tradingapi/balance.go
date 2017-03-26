package tradingapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Balances map[string]float64

// Poloniex trading API implementation of returnBalances command.
//
// API Doc:
// Returns all of your available balances.
//
// Sample output:
//
//  {
//    "BTC": "0.59098578",
//    "LTC": "3.31117268", ...
//  }
func (client *TradingClient) GetBalances() (Balances, error) {

	resp, err := client.do("returnBalances")
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	var res = make(map[string]string)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	resConv := make(Balances)

	for key, value := range res {

		if res, err := strconv.ParseFloat(value, 64); err != nil {
			return nil, fmt.Errorf("parsefloat: %v", err)
		} else {
			resConv[key] = res
		}
	}
	return resConv, nil
}
