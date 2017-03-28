package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
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

	postParameters := url.Values{}
	postParameters.Add("command", "returnBalances")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("do: %v", err)
	}

	res := make(Balances)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}

func (b *Balances) UnmarshalJSON(data []byte) error {

	res := make(map[string]string)

	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("json unmarshal: %v", err)
	}

	*b = make(Balances)
	for key, value := range res {

		if res, err := strconv.ParseFloat(value, 64); err != nil {
			return fmt.Errorf("parsefloat: %v", err)
		} else {
			(*b)[key] = res
		}
	}

	return nil
}
