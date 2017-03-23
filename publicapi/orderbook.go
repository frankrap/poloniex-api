package publicapi

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
	IsFrozen bool
	Seq      uint64
}

type Order struct {
	Rate     float64
	Quantity float64
}

func (o *OrderBook) UnmarshalJSON(data []byte) error {

	type Alias OrderBook
	aux := struct {
		IsFrozen string
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("unmarshal aux: %v", err)
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
		return err
	}

	if got, want := len(tmp), 2; got != want {
		return fmt.Errorf("wrong number of fields in Order: %d != %d",
			got, want)
	}

	val, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		return fmt.Errorf("parsefloat: %v", err)
	}
	o.Rate = val

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
