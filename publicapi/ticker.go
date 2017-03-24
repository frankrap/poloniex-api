package publicapi

import (
	"encoding/json"
	"fmt"
)

type AllTicks map[string]Tick

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

func (t *Tick) UnmarshalJSON(data []byte) error {

	type alias Tick
	aux := struct {
		IsFrozen string
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal aux: %v", err)
	}
	fmt.Println(aux.IsFrozen)
	if aux.IsFrozen != "0" {
		t.IsFrozen = true
	} else {
		t.IsFrozen = false
	}

	return nil
}

func (client *PublicClient) GetTicker() (AllTicks, error) {

	params := map[string]string{
		"command": "returnTicker",
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := make(AllTicks)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}
