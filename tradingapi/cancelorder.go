package tradingapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type CanceledOrder struct {
	Success bool    `json:"success"`
	Amount  float64 `json:"amount,string"`
	Message string  `json:"message"`
}

// Poloniex trading API implementation of cancelOrder command.
//
// API Doc:
// Cancels an order you have placed in a given market. Required POST parameter is "orderNumber".
//
// Sample output:
//
//  {
//    "success": 1
//    "amount": "0.1"
//    "message": "Order #258128814946 canceled."
//  }
func (client *Client) CancelOrder(orderNumber int64) (*CanceledOrder, error) {

	postParameters := url.Values{}
	postParameters.Add("command", "cancelOrder")
	postParameters.Add("orderNumber", strconv.Itoa(int(orderNumber)))

	resp, err := client.do(postParameters)
	if err != nil {
		return nil, fmt.Errorf("TradingClient.do: %v", err)
	}

	res := CanceledOrder{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return &res, nil
}

func (c *CanceledOrder) UnmarshalJSON(data []byte) error {

	type alias CanceledOrder
	aux := struct {
		Success int `json:"success"`
		*alias
	}{
		alias: (*alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("json.Unmarshal: %v", err)
	}

	if aux.Success != 1 {
		c.Success = false
	} else {
		c.Success = true
	}

	return nil
}
