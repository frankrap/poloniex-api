package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DayVolumes struct {
	DayVolumes      map[string]*DayVolume
	PrimaryCurrency map[string]float64
}

type DayVolume map[string]float64

// Poloniex public API implementation of return24Volume command.
//
// API Doc:
// Returns the 24-hour volume for all markets, plus totals for primary currencies.
//
// Call: https://poloniex.com/public?command=return24hVolume
//
// Sample output:
//
//  {
//    "BTC_BBR": {
//      "BTC": "22.86850902",
//      "BBR": "89588.33257239"
//    },
//    "BTC_BCN": {
//      "BTC": "4.10236135",
//      "BCN": "89984007.28797383"
//    }, ...
//    "totalBTC": "119908.58082298",
//    "totalETH": "13207.76114161",
//    "totalUSDT": "23533800.94795309",
//    "totalXMR": "3675.72265894",
//    "totalXUSD": "0.00000000"
//  }
func (client *Client) GetDayVolumes() (*DayVolumes, error) {

	params := map[string]string{
		"command": "return24hVolume",
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := DayVolumes{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

func (dv *DayVolumes) UnmarshalJSON(data []byte) error {

	adv := make(map[string]interface{})

	if err := json.Unmarshal(data, &adv); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	dv.DayVolumes = make(map[string]*DayVolume)
	dv.PrimaryCurrency = make(map[string]float64)

	for key, value := range adv {

		switch value := value.(type) {

		case map[string]interface{}:
			if res, err := convertToDayVolume(value); err != nil {
				return fmt.Errorf("convertToDayVolume: %v", err)
			} else {
				dv.DayVolumes[key] = res
			}

		case string:
			if res, err := strconv.ParseFloat(value, 64); err != nil {
				return fmt.Errorf("strconv.ParseFloat: %v", err)
			} else {
				dv.PrimaryCurrency[key] = res
			}

		default:
			return fmt.Errorf("Type error %v", value)
		}
	}

	return nil
}

func convertToDayVolume(value map[string]interface{}) (*DayVolume, error) {

	dv := make(DayVolume)
	for k, v := range value {

		if v, ok := v.(string); ok {

			if val, err := strconv.ParseFloat(v, 64); err != nil {
				return nil, fmt.Errorf("strconv.ParseFloat : %v", err)
			} else {
				dv[k] = val
			}

		} else {
			return nil, fmt.Errorf("Type error: %v", v)
		}
	}
	return &dv, nil
}
