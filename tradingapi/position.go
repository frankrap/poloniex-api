package tradingapi

import (
	"fmt"
	"net/url"
	"encoding/json"
)

type Positions map[string]Position

/*
{"amount":"40.94717831","total":"-0.09671314",""basePrice":"0.00236190","liquidationPrice":-1,"pl":"-0.00058655", "lendingFees":"-0.00000038","type":"long"}
*/
type Position struct {
	Amount           float64 `json:"amount,string"`
	Total            float64 `json:"total,string"`
	BasePrice        float64 `json:"basePrice,string"`
	LiquidationPrice float64 `json:"liquidationPrice"`
	PL               float64 `json:"pl,string"`
	LendingFees      float64 `json:"lendingFees,string"`
	Type             string  `json:"type"`
}

func (client *Client) GetMarginPosition(currencyPair string) (*Positions, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "getMarginPosition")
	postParameters.Add("currencyPair", currencyPair)

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := Positions{}

	if string(resp) == "" {
		return &res, nil
	}
	fmt.Println(string(resp))

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
