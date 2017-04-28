package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type TradableBalances map[string]TradableBalance

type TradableBalance map[string]float64

// Poloniex trading API implementation of returnTradableBalances command.
//
// API Doc:
// Returns your current tradable balances for each currency in each market for which margin
// trading is enabled. Please note that these balances may vary continually with market conditions.
//
// Sample output:
//
//  {
//    "BTC_DASH": {
//      "BTC": "8.50274777",
//      "DASH": "654.05752077"
//    },
//    "BTC_LTC": {
//      "BTC": "8.50274777",
//      "LTC": "1214.67825290"
//    }, ...
//  }
func (client *Client) GetTradableBalances() (TradableBalances, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnTradableBalances")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(TradableBalances)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (t *TradableBalances) UnmarshalJSON(data []byte) error {

	res := make(map[string]map[string]string)

	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("json.Umarshal: %v", err)
	}

	*t = make(TradableBalances)

	for currencyPair, tradableBalance := range res {

		(*t)[currencyPair] = make(TradableBalance)
		for cur, val := range tradableBalance {

			if r, err := strconv.ParseFloat(val, 64); err != nil {
				return fmt.Errorf("strconv.ParseFloat: %v", err)
			} else {
				(*t)[currencyPair][cur] = r
			}
		}
	}

	return nil
}
