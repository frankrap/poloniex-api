package publicapi

import (
	"encoding/json"
	"fmt"
)

type Currencies map[string]Currency

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
		return fmt.Errorf("unmarshal aux: %v", err)
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

func (client *PublicClient) GetCurrencies() (Currencies, error) {

	params := map[string]string{
		"command": "returnCurrencies",
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	var res = make(Currencies)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}
