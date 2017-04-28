package publicapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LoanOrders struct {
	Offers  []*LoanOrder `json:"offers"`
	Demands []*LoanOrder `json:"demands"`
}

type LoanOrder struct {
	Rate     float64 `json:"rate,string"`
	Amount   float64 `json:"amount,string"`
	RangeMin int     `json:"rangeMin"`
	RangeMax int     `json:"rangeMax"`
}

// Poloniex public API implementation of returnLoanOrders command.
//
// API Doc:
// Returns the list of loan offers and demands for a given currency,
// specified by the "currency" GET parameter.
//
// Call: https://poloniex.com/public?command=returnLoanOrders&currency=BTC
//
// Sample output:
//
//  {
//    "offers": [
//      {
//        "rate": "0.00288800",
//        "amount": "0.49414692",
//        "rangeMin": 2,
//        "rangeMax": 2
//      },
//      {
//        "rate": "0.00288900",
//        "amount": "0.05031184",
//        "rangeMin": 2,
//        "rangeMax": 2
//      }, ...
//    "demands": [
//      {
//        "rate": "0.00200000",
//        "amount": "0.32648833",
//        "rangeMin": 2,
//        "rangeMax": 2
//      },
//      {
//        "rate": "0.00120100",
//        "amount": "2.49999988",
//        "rangeMin": 2,
//        "rangeMax": 2
//      }, ...
//    ]
//  }
func (client *Client) GetLoanOrders(currency string) (*LoanOrders, error) {

	params := map[string]string{
		"command":  "returnLoanOrders",
		"currency": strings.ToUpper(currency),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := LoanOrders{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}
