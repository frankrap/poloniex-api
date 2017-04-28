package publicapi

import (
	"encoding/json"
	"fmt"
)

type Currencies map[string]*Currency

type Currency struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	TxFee          float64 `json:"txFee,string"`
	MinConf        int     `json:"minConf"`
	DepositAddress string  `json:"depositAddress"`
	Disabled       bool
	Delisted       bool
	Frozen         bool
}

// Poloniex public API implementation of returnCurrencies command.
//
// API Doc:
// Returns information about currencies.
//
// Call: https://poloniex.com/public?command=returnCurrencies
//
// Sample output:
//
//  {
//    "1CR": {
//      "id": 1,
//      "name": "1CRedit",
//      "txFee": "0.01000000",
//      "minConf": 3,
//      "depositAddress": null,
//      "disabled": 0,
//      "delisted": 1,
//      "frozen": 0
//    }, ...
//  }
func (client *Client) GetCurrencies() (Currencies, error) {

	params := map[string]string{
		"command": "returnCurrencies",
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	var res = make(Currencies)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (c *Currency) UnmarshalJSON(data []byte) error {

	type alias Currency
	aux := struct {
		Disabled int `json:"disabled"`
		Delisted int `json:"delisted"`
		Frozen   int `json:"frozen"`
		*alias
	}{
		alias: (*alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if aux.Disabled != 0 {
		c.Disabled = true
	} else {
		c.Disabled = false
	}

	if aux.Delisted != 0 {
		c.Delisted = true
	} else {
		c.Delisted = false
	}

	if aux.Frozen != 0 {
		c.Frozen = true
	} else {
		c.Frozen = false
	}

	return nil
}

func (collection Currencies) Filter(f func(*Currency) bool) Currencies {

	res := make(Currencies)

	for k, currency := range collection {
		if f(currency) {
			res[k] = currency
		}
	}
	return res
}
