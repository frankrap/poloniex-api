package publicapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type OrderBooks map[string]*OrderBook

type OrderBook struct {
	Asks     []*Order `json:"asks"`
	Bids     []*Order `json:"bids"`
	IsFrozen bool
	Seq      int64
}

type Order struct {
	Rate     float64
	Quantity float64
}

// Poloniex public API implementation of returnOrderBook command.
//
// API Doc:
// Returns the order book for a given market, as well as a sequence number for use with
// the Push API and an indicator specifying whether the market is frozen. You may set
// currencyPair to "all" to get the order books of all markets.
//
// Call: https://poloniex.com/public?command=returnOrderBook&currencyPair=BTC_NXT&depth=10
//
// Sample output:
//
//  {
//    "asks": [
//      [
//        "0.00001315",
//        36937.09233522
//      ],
//      [
//        "0.00001332",
//        8365.874
//      ], ...
//    ],
//    bids": [
//      [
//        "0.00001311",
//        6006.00485372
//      ],
//      [
//        "0.00001309",
//        6602.96320483
//      ], ...
//    ]
//    "isFrozen": "0",
//    "seq": 28233022
//  }
func (client *Client) GetOrderBook(currencyPair string, depth int) (*OrderBook, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": strings.ToUpper(currencyPair),
		"depth":        strconv.Itoa(depth),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := OrderBook{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

// GetOrderBooks returns the order books for all markets (currencyPair to "all")
//
// Call: https://poloniex.com/public?command=returnOrderBook&currencyPair=ALL&depth=10
//
// Sample output:
//
//  {
//    "BTC_AMP": {
//      "asks": [
//        [
//          "0.00006371",
//          41.23554691
//        ],
//        [
//          "0.00006386",
//          4071.27563735
//        ], ...
//      ]
//      "bids": [
//        [
//          "0.00006356",
//          54.15144965
//        ],
//        [
//          "0.00006353",
//          23811.21533134
//        ], ...
//      ]
//      "isFrozen": "0",
//      "seq": 30446838
//    }, ...
//  }
func (client *Client) GetOrderBooks(depth int) (OrderBooks, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": "all",
		"depth":        strconv.Itoa(depth),
	}

	resp, err := client.do(params)
	if err != nil {
		return nil, fmt.Errorf("PublicClient.do: %v", err)
	}

	res := make(OrderBooks)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return res, nil
}

func (o *OrderBook) UnmarshalJSON(data []byte) error {

	type alias OrderBook
	aux := struct {
		IsFrozen string
		*alias
	}{
		alias: (*alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if aux.IsFrozen != "0" {
		o.IsFrozen = false
	} else {
		o.IsFrozen = true
	}

	return nil
}

func (o *Order) UnmarshalJSON(data []byte) error {

	var rateStr string
	tmp := []interface{}{&rateStr, &o.Quantity}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if got, want := len(tmp), 2; got != want {
		return fmt.Errorf("wrong number of fields in Order: %d != %d",
			got, want)
	}

	if val, err := strconv.ParseFloat(rateStr, 64); err != nil {
		return fmt.Errorf("strconv.ParseFloat: %v", err)
	} else {
		o.Rate = val
	}

	return nil
}
