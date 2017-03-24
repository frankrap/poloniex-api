package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ChartData []CandleStick

type CandleStick struct {
	Date            int64   `json:"date"` // Unix timestamp
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	Close           float64 `json:"close"`
	Volume          float64 `json:"volume"`
	QuoteVolume     float64 `json:"quoteVolume"`
	WeighedtAverage float64 `json:"weightedAverage"`
}

func (client *PublicClient) GetChartData(currencyPair string, start, end time.Time, period int) (ChartData, error) {

	switch period { // Valid period only
	case 300: // 5min
	case 900: // 15min
	case 1800: // 30min
	case 7200: // 2h
	case 14400: // 4h
	case 86400: // 1d
	default:
		return nil, fmt.Errorf("wrong period: %d", period)
	}

	params := map[string]string{
		"command":      "returnChartData",
		"currencyPair": strings.ToUpper(currencyPair),
		"start":        strconv.Itoa(int(start.Unix())),
		"end":          strconv.Itoa(int(end.Unix())),
		"period":       strconv.Itoa(period),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	var res = make(ChartData, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}
