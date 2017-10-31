package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	poloniex "github.com/joemocquant/poloniex-api"
)

// Poloniex trading API implementation of moveOrder command.
//
// API Doc:
// Cancels an order and places a new one of the same type in a single atomic transaction,
// meaning either both operations will succeed or both will fail. Required POST parameters
// are "orderNumber" and "rate"; you may optionally specify "amount" if you wish to change
// the amount of the new order. "postOnly" or "immediateOrCancel" may be specified for exchange
// orders, but will have no effect on margin orders.
//
// Sample output:
//
//  {
//    "success": 1,
//    "orderNumber": "239574176",
//    "resultingTrades": {
//      "BTC_BTS": [
//        {
//          "amount": "338.8732",
//          "date": "2014-10-18 23:03:21",
//          "rate": "0.00000173",
//          "total": "0.00058625",
//          "tradeID": "16164",
//          "type": "buy"
//        }, ...
//      ]
//    }
//  }
type MovedOrder struct {
	Success         bool                                  `json:"success"`
	OrderNumber     int64                                 `json:"orderNumber,string"`
	ResultingTrades map[string][]*poloniex.ResultingTrade `json:"resultingTrades"`
}

func (client *Client) MoveOrderPostOnly(orderNumber int64, rate, amount float64) (*MovedOrder, error) {
	return client.moveOrder(orderNumber, rate, amount, "postOnly")
}

func (client *Client) MoveOrderImmediateOrCancel(orderNumber int64, rate, amount float64) (*MovedOrder, error) {
	return client.moveOrder(orderNumber, rate, amount, "immediateOrCancel")
}

func (client *Client) MoveOrder(orderNumber int64, rate, amount float64) (*MovedOrder, error) {
	return client.moveOrder(orderNumber, rate, amount, "")
}

func (client *Client) moveOrder(orderNumber int64, rate, amount float64, option string) (*MovedOrder, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "moveOrder")
	postParameters.Add("orderNumber", strconv.Itoa(int(orderNumber)))
	postParameters.Add("rate", strconv.FormatFloat(rate, 'f', -1, 64))
	postParameters.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	if option != "" {
		postParameters.Add(option, "1")
	}

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := MovedOrder{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

func (m *MovedOrder) UnmarshalJSON(data []byte) error {

	type alias MovedOrder
	aux := struct {
		Success int `json:"success"`
		*alias
	}{
		alias: (*alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if aux.Success != 1 {
		m.Success = false
	} else {
		m.Success = true
	}

	return nil
}
