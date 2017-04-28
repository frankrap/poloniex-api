package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Poloniex trading API implementation of generateNewAddress command.
//
// API Doc:
// Generates a new deposit address for the currency specified by the "currency" POST parameter.
//
// Sample output:
//
//  {
//    "success": 1,
//    "response": "CKXbbs8FAVbtEa397gJHSutmrdrBrhUMxe"
//  }
func (client *Client) GenerateNewAddress(currency string) (string, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "generateNewAddress")
	postParameters.Add("currency", currency)

	resp, err := client.do(postParameters)
	if err != nil {
		return "", fmt.Errorf("TradingClient.do: %v", err)
	}

	type Result struct {
		Success  int    `json:"success"`
		Response string `json:"response"`
	}

	res := Result{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return "", fmt.Errorf("json.Unmarshal: %v", err)
	}

	if res.Success != 1 {
		return "", fmt.Errorf("Error response: %s", res.Response)
	}

	return res.Response, nil
}
