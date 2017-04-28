package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type FeeInfo struct {
	MakerFee        float64 `json:"makerFee,string"`
	TakerFee        float64 `json:"takerFee,string"`
	ThirtyDayVolume float64 `json:"thirtyDayVolume,string"`
	NextTier        float64 `json:"nextTier,string"`
}

// Poloniex trading API implementation of returnFeeInfo command.
//
// API Doc:
// If you are enrolled in the maker-taker fee schedule, returns your current trading fees
// and trailing 30-day volume in BTC. This information is updated once every 24 hours.
//
// Sample output:
//
//  {
//    "makerFee": "0.00140000",
//    "takerFee": "0.00240000",
//    "thirtyDayVolume": "612.00248891",
//    "nextTier": "1200.00000000"
//  }
func (client *Client) GetFeeInfo() (*FeeInfo, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "returnFeeInfo")

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := FeeInfo{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
