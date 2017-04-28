package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ChartData []*CandleStick

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

// Poloniex public API implementation of returnChartData command.
//
// API Doc:
// Returns candlestick chart data. Required GET parameters are "currencyPair", "period"
// (candlestick period in seconds; valid values are 300, 900, 1800, 7200, 14400, and 86400),
// "start", and "end". "Start" and "end" are given in UNIX timestamp format and used to specify
// the date range for the data returned.
//
// Call: https://poloniex.com/public?command=returnChartData&currencyPair=BTC_XMR&start=1405699200&end=9999999999&period=14400
//
// Sample output:
//
//  [
//    {
//      "date": 1405699200,
//      "high": 0.0045388,
//      "low": 0.00403001,
//      "open": 0.00404545,
//      "close": 0.00427592,
//      "volume": 44.11655644,
//      "quoteVolume": 10259.29079097,
//      "weightedAverage": 0.00430015
//    }, ...
//  ]
func (client *Client) GetChartData(currencyPair string, start, end time.Time, period int) (ChartData, error) {

	switch period { // Valid period only
	case 300: // 5min
	case 900: // 15min
	case 1800: // 30min
	case 7200: // 2h
	case 14400: // 4h
	case 86400: // 1d
	default:
		return nil, fmt.Errorf("Wrong period parameter: %d", period)
	}

	params := map[string]string{
		"command":      "returnChartData",
		"currencyPair": strings.ToUpper(currencyPair),
		"start":        strconv.Itoa(int(start.Unix())),
		"end":          strconv.Itoa(int(end.Unix())),
		"period":       strconv.Itoa(period),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	var res = make(ChartData, 200)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}
