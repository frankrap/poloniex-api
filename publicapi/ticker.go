package publicapi

import (
	"encoding/json"
	"fmt"
)

type Ticks map[string]*Tick

type Tick struct {
	Id            int     `json:"id"`
	Last          float64 `json:"last,string"`
	LowestAsk     float64 `json:"lowestAsk,string"`
	HighestBid    float64 `json:"highestBid,string"`
	PercentChange float64 `json:"percentChange,string"`
	BaseVolume    float64 `json:"baseVolume,string"`
	QuoteVolume   float64 `json:"quoteVolume,string"`
	IsFrozen      bool
	High24hr      float64 `json:"high24hr,string"`
	Low24hr       float64 `json:"low24hr,string"`
}

// Poloniex public API implementation of returnTicker command.
//
// API Doc:
// Returns the ticker for all markets.
//
// Call: https://poloniex.com/public?command=returnTicker
//
// Sample output:
//
//  {
//    "BTC_BBR": {
//      "id": 6,
//      "last": "0.00024306",
//      "lowestAsk": "0.00024306",
//      "highestBid": "0.00024305",
//      "percentChange": "-0.10662697",
//      "baseVolume": "21.86898934",
//      "quoteVolume": "85944.61508131",
//      "isFrozen": "0",
//      "high24hr": "0.00027359",
//      "low24hr": "0.00023653"
//    }, ...
//  }
func (client *Client) GetTickers() (Ticks, error) {

	params := map[string]string{
		"command": "returnTicker",
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := make(Ticks)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (t *Tick) UnmarshalJSON(data []byte) error {

	type alias Tick
	aux := struct {
		IsFrozen string
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if aux.IsFrozen != "0" {
		t.IsFrozen = true
	} else {
		t.IsFrozen = false
	}

	return nil
}
