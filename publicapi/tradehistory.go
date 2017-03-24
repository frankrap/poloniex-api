package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TradeHistory []Trade

type Trade struct {
	GlobalTradeId int     `json:"globalTradeID"`
	TradeId       int     `json:"tradeID"`
	Date          int64   // Unix timestamp
	TypeOrder     string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total         float64 `json:"total,string"`
}

func (t *Trade) UnmarshalJSON(data []byte) error {

	type alias Trade
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal aux: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("timestamp conversion: %v", err)
	} else {
		t.Date = int64(timestamp.Unix())
	}

	return nil
}

func (client *PublicClient) GetPast200TradeHistory(currencyPair string) (TradeHistory, error) {

	params := map[string]string{
		"command":      "returnTradeHistory",
		"currencyPair": strings.ToUpper(currencyPair),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := make(TradeHistory, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}

func (client *PublicClient) GetTradeHistory(currencyPair string, start, end time.Time) (TradeHistory, error) {

	params := map[string]string{
		"command":      "returnTradeHistory",
		"currencyPair": strings.ToUpper(currencyPair),
		"start":        strconv.Itoa(int(start.Unix())),
		"end":          strconv.Itoa(int(end.Unix())),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	var res = make(TradeHistory, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}
