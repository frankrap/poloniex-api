// Poloniex commons
package poloniex

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func PrettyPrintJson(msg interface{}) {

	jsonstr, err := json.MarshalIndent(msg, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(jsonstr))
}

type ResultingTrade struct {
	Amount    float64 `json:"amount,string"`
	Date      int64   // Unix timestamp
	Rate      float64 `json:"rate,string"`
	Total     float64 `json:"total,string"`
	TradeId   int64   `json:"tradeID,string"`
	TypeOrder string  `json:"type"`
}

func (r *ResultingTrade) UnmarshalJSON(data []byte) error {

	type alias ResultingTrade
	aux := struct {
		Date string `json:"Date"`
		*alias
	}{
		alias: (*alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if timestamp, err := time.Parse("2006-01-02 15:04:05", aux.Date); err != nil {
		return fmt.Errorf("time.Parse: %v", err)
	} else {
		r.Date = int64(timestamp.Unix())
	}

	return nil
}
