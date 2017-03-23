package poloniex

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type AllOrderBooks map[string]OrderBook

type OrderBook struct {
	Asks     []Order `json:"asks"`
	Bids     []Order `json:"bids"`
	IsFrozen string  `json:"isFrozen"`
	Seq      float64
}

type Order struct {
	Rate     string
	Quantity float64
}

func (o *Order) UnmarshalJSON(buf []byte) error {

	tmp := []interface{}{&o.Rate, &o.Quantity}

	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}

	if got, want := len(tmp), 2; got != want {
		return fmt.Errorf("wrong number of fields in Order: %d != %d", got, want)
	}

	return nil
}

func (client *PublicClient) GetAllOrderBooks(depth int) (AllOrderBooks, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": "all",
		"depth":        strconv.Itoa(depth),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	var res = make(AllOrderBooks)

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return res, nil
}

func (client *PublicClient) GetOrderBook(currencyPair string, depth int) (*OrderBook, error) {

	params := map[string]string{
		"command":      "returnOrderBook",
		"currencyPair": strings.ToUpper(currencyPair),
		"depth":        strconv.Itoa(depth),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	var res OrderBook

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}
