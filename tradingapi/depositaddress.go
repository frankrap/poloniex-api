package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type DepositAddresses map[string]string

// Poloniex trading API implementation of returnDepositAddresses command.
//
// API Doc:
// Returns all of your deposit addresses.
//
// Sample output:
//
// {
//     "BTC": "19YqztHmspv2egyD6jQM3yn81x5t5krVdJ", ...
// }
func (client *Client) GetDepositAddresses() (DepositAddresses, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnDepositAddresses")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := make(DepositAddresses)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}
